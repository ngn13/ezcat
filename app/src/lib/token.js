function tokenSet(t) {
  localStorage.setItem("token", t);
}

function tokenGet() {
  return localStorage.getItem("token");
}

function tokenClear() {
  localStorage.removeItem("token");
}

export { tokenSet, tokenGet, tokenClear };
