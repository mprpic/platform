import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Formats a timestamp string (RFC3339) into a human-readable format
 * @param timestamp - RFC3339 timestamp string
 * @returns Formatted time string (e.g., "2:30 PM" or "Jan 15, 2:30 PM")
 */
export function formatMessageTimestamp(timestamp: string): string {
  try {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    // If less than 1 minute ago, show "Just now"
    if (diffMins < 1) {
      return "Just now";
    }

    // If less than 60 minutes ago, show "X min ago"
    if (diffMins < 60) {
      return `${diffMins} min ago`;
    }

    // If less than 24 hours ago, show "X hour(s) ago"
    if (diffHours < 24) {
      return `${diffHours} hour${diffHours === 1 ? "" : "s"} ago`;
    }

    // If less than 7 days ago, show "X day(s) ago"
    if (diffDays < 7) {
      return `${diffDays} day${diffDays === 1 ? "" : "s"} ago`;
    }

    // Otherwise show date and time
    const isToday = date.toDateString() === now.toDateString();
    const isThisYear = date.getFullYear() === now.getFullYear();

    const timeStr = date.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    });

    if (isToday) {
      return timeStr;
    }

    const dateStr = date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      ...(isThisYear ? {} : { year: "numeric" }),
    });

    return `${dateStr}, ${timeStr}`;
  } catch {
    return timestamp;
  }
}
