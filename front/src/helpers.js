import endpoints from "./endpoints";
/**
 * @param {string} url
 */
export async function fetchOrFail(url) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 2000);

  const response = await fetch(url, { signal: controller.signal });
  if (response.ok) {
    clearTimeout(timeoutId);
    return await response.json();
  }

  console.log(response.error);
}

/**
 * sends command to server
 * @param {string} url
 * @param {string} command
 * @param {object} data
 */
export async function sendCommand(url, command, data) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 2000);
  const request = {
    command: command,
    data: data,
  };
  const response = await fetch(url, {
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
/**
 * @param {string} directory
 * @return {object}
 */
export async function getFolders(directory) {
  const request = {
    command: "list",
    data: {
      playlist: directory,
    },
  };
  const response = await fetch(endpoints.folders, {
    method: "POST",
    headers: {
      "Content-Type": "application/json;charset=utf-8",
    },
    body: JSON.stringify(request),
  });
  if (response.ok) {
    const json = await response.json();
    return [...json.Folders, ...json.Files];
  }
  console.log(response.error);
  return [];
}

/**
 * returns time in h:m:s format
 * @param {number} time
 * @return {string}
 */
export function humanizeTime(time) {
  function spanZero(time) {
    return time < 10 ? "0" + time : time;
  }
  let second = Math.floor(Number(time) % 60);
  if (time < 3600) {
    const minute = Math.floor(Number(time) / 60);
    second = spanZero(second);
    return minute + ":" + second;
  }
  const hour = Math.floor(Number(time) / 3600);
  let minute = Math.floor((time % 3600) / 60);
  second = spanZero(second);
  minute = spanZero(minute);
  return hour + ":" + minute + ":" + second;
}
/** toggles media controller visibility in mobile */
export function toggleMediaController() {
  document.getElementById("mediaController").classList.toggle("z-10");
  document.getElementById("queue").classList.toggle("z-10");
}

/**  toggles media controller visibility in mobile */
export function toggleFolders() {
  document.getElementById("mediaController").classList.toggle("z-10");
  document.getElementById("sidebar").classList.toggle("z-10");
}
