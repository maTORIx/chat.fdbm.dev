import { ChangeEvent } from "react";

export function withEventValue(f: Function): (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void {
  return (e) => f(e.target.value);
}

export function generateRandomString(len: Number): String {
  const charList =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array.from(Array(len))
    .map(() => charList[Math.floor(Math.random() * charList.length)])
    .join("");
}

export function generateRandomColorFromText(str: string): string {
  if (str.length < 1) {
    return "#ffffff";
  }
  while (str.length < 6) {
    str += str
  }
  let result = [0, 0, 0, 0, 0, 0]
  str.split("").forEach((v, idx) => result[idx % 6] += v.charCodeAt(0))
  return "#" + result.map((v) => (v % 16).toString(16)).join("")
}