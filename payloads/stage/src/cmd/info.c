#ifdef _WIN32
// clang-format off
#include <winsock2.h>
#include <windows.h>
// clang-format on
#else
#include <sys/socket.h>
#include <sys/utsname.h>
#endif

#include <limits.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "../cmd.h"
#include "../util.h"

bool cmd_info_handler(agent_t *agent, packet_t *packet) {
  ssize_t infolen = 0;

#ifdef _WIN32
#define UNLEN 256
#define CNLEN 15

  char  hostname[CNLEN + 1];
  char  username[UNLEN + 1];
  DWORD username_len = UNLEN + 1;
  HKEY  hkey;
  DWORD keysize = 0;
  DWORD keytype = REG_SZ;

  bzero(hostname, CNLEN + 1);
  bzero(username, UNLEN + 1);

  if (gethostname(hostname, CNLEN + 1) != 0) {
    debug("failed to get hostname");
    cmd_failure(agent, "failed to get the hostname", 0);
    return false;
  }

  if (GetUserName(username, &username_len) == 0) {
    debug("failed to get username");
    cmd_failure(agent, "failed to get the username", 0);
    return false;
  }

  if (RegOpenKeyExW(HKEY_LOCAL_MACHINE, L"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", 0, KEY_READ, &hkey) !=
      ERROR_SUCCESS) {
    debug("failed to get version info");
    cmd_failure(agent, "failed to get version info", 0);
    return false;
  }

  if (RegQueryValueExW(hkey, L"ProductName", NULL, NULL, NULL, &keysize) != ERROR_SUCCESS) {
    debug("failed to query key size");
    cmd_failure(agent, "failed to query key size", 0);
    RegCloseKey(hkey);
    return false;
  }

  char version[keysize + 1];
  bzero(version, keysize + 1);
  RegCloseKey(hkey);

  if (RegGetValueA(HKEY_LOCAL_MACHINE,
          "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion",
          "ProductName",
          RRF_RT_REG_SZ,
          NULL,
          version,
          &keysize) != ERROR_SUCCESS) {
    debug("failed to read the key");
    cmd_failure(agent, "failed to read the key", 0);
    RegCloseKey(hkey);
    return false;
  }

  // 10 is for the PID, 5 is for extra space (@ and the space)
  infolen = strlen(hostname) + strlen(username) + keysize + 1 + 10 + 5;

  char info[infolen];
  bzero(info, infolen);

  snprintf(info, infolen, "%s@%s@%s@%d", username, hostname, version, getpid());
#else
  char           hostname[HOST_NAME_MAX + 1];
  char           username[LOGIN_NAME_MAX + 1];
  char          *distro = get_distro();
  struct utsname udata;

  bzero(hostname, HOST_NAME_MAX + 1);
  bzero(username, LOGIN_NAME_MAX + 1);

  if (gethostname(hostname, HOST_NAME_MAX + 1) != 0) {
    debug("failed to get hostname");
    cmd_failure(agent, "failed to get the hostname", 0);
    return false;
  }

  if (getlogin_r(username, LOGIN_NAME_MAX + 1) != 0) {
    debug("failed to get username");
    cmd_failure(agent, "failed to get the username", 0);
    return false;
  }

  if (uname(&udata) < 0) {
    debug("failed to get uname");
    cmd_failure(agent, "failed to get uname", 0);
    return false;
  }

  // 10 is for the PID, 5 is for extra space (@ and the space)
  infolen = strlen(username) + strlen(hostname) + strlen(udata.sysname) + strlen(udata.release) + 10 + 5;

  if (NULL != distro)
    infolen += strlen(distro);

  char info[infolen];
  bzero(info, infolen);

  if (NULL == distro)
    snprintf(info, infolen, "%s@%s@%s %s@%d", username, hostname, udata.sysname, udata.release, getpid());
  else
    snprintf(info, infolen, "%s@%s@%s %s %s@%d", username, hostname, distro, udata.sysname, udata.release, getpid());
#endif

  debug("sending %s", info);

  if (!cmd_success(agent, info, 0)) {
    debug("failed to send the info command result");
    return false;
  }

  return true;
}
