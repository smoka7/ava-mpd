import endpoints from "./endpoints"
/**
 * @param {string} url
 */
export async function fetchOrFail(url) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 2000);

  let response = await fetch(url, { signal: controller.signal });
  if (response.ok) {
    clearTimeout(timeoutId);
    return await response.json();
  }

  console.log(response.error);
}

export async function sendCommand(url, command, data) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 2000);
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
    clearTimeout(timeoutId);
    return;
  }
  console.log(response.error);
}
export async function getFolders(directory) {
  let request = {
    command: "list",
    data: {
      playlist: directory,
    },
  };
  let response = await fetch(endpoints.folders, {
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
    body: JSON.stringify(request),
  });
  if (response.ok) {
    let json = await response.json();
    return [...json.Folders, ...json.Files];
  }
  console.log(response.error);
  return [];
}

// returns time in h:m:s format
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
// toggles media controller visibility in mobile
export function toggleMediaController() {
  document.getElementById("mediaController").classList.toggle("z-10");
  document.getElementById("queue").classList.toggle("z-10");
}
// toggles media controller visibility in mobile
export function toggleFolders() {
  document.getElementById("mediaController").classList.toggle("z-10");
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
