export async function sendCommand(url, command, data) {
  let request = {
    command: command,
    data: data,
  };
  let response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
    body: JSON.stringify(request),
  });
  if (response.ok) {
    return;
  }
  console.log(response.error);
}
//returns time in h:m:s format
export function humanizeTime(time) {
  function spanZero(time) {
    return time < 10 ? "0" + time : time;
  }
  let second = Math.floor(Number(time) % 60);
  if (time < 3600) {
    let minute = Math.floor(Number(time) / 60);
    second = spanZero(second);
    return minute + ":" + second;
  }
  let hour = Math.floor(Number(time) / 3600);
  let minute = Math.floor((time % 3600) / 60);
  second = spanZero(second);
  minute = spanZero(minute);
  return hour + ":" + minute + ":" + second;
}
//toggles media controller visibility in mobile
export function toggleMediaController() {
  document.getElementById('mediaController').classList.toggle('z-10');
  document.getElementById("queue").classList.toggle("z-10");
}
//toggles media controller visibility in mobile
export function toggleFolders() {
  document.getElementById('mediaController').classList.toggle('z-10');
  document.getElementById("sidebar").classList.toggle("z-10");
}

export function toggleDarkMode() {
  if (localStorage.theme === "dark") {
    document.documentElement.classList.remove("dark");
    localStorage.theme = "light";
    this.theme = "light";
  } else {
    document.documentElement.classList.add("dark");
    localStorage.theme = "dark";
    this.theme = "dark";
  }
}
