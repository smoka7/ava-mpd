import type endpoints from "./endpoints";

/**
 * @param {string} url
 */
export async function fetchOrFail<T>(url: string): Promise<T> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 2000);

  const response = await fetch(url, { signal: controller.signal });
  if (response.ok) {
    clearTimeout(timeoutId);
    return await response.json();
  }
  console.log(response);
  return await response.json();
}

/**
 * sends command to server
 * @param {string} url
 * @param {string} command
 * @param {object} data
 */
export async function sendCommand(url: endpoints, command: string, data?: any) {
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
    try {
      const res = await response.json();
      return res;
    } catch (error) {
      return;
    }
  }
}

/**
 * returns time in h:m:s format
 * @param {number} time
 * @return {string}
 */
export function humanizeTime(time: number): string {
  function spanZero(time: number): string {
    return time < 10 ? "0" + time : time.toString();
  }
  const second = Math.floor(Number(time) % 60);
  if (time <= 60) {
    return 0 + ":" + second.toString();
  }
  const minute = Math.floor(Number(time) / 60);
  if (time < 3600) {
    return minute + ":" + spanZero(second);
  }
  const hour = Math.floor(Number(time) / 3600);
  return hour + ":" + spanZero(minute) + ":" + spanZero(second);
}

/** toggles media controller visibility in mobile */
export function toggleMediaController() {
  const cl = document.getElementById("mediaController");
  const queue = document.getElementById("queue");
  if (cl === null || queue == null) {
    return;
  }
  cl.classList.toggle("z-10");
  queue.classList.toggle("z-10");
}
