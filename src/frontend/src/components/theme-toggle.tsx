"use client";

import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";

const themeOrder = ["system", "light", "dark"] as const;

const themeIcons = {
  system: Monitor,
  light: Sun,
  dark: Moon,
} as const;

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  function cycleTheme() {
    const currentIndex = themeOrder.indexOf(theme as (typeof themeOrder)[number]);
    const nextIndex = (currentIndex + 1) % themeOrder.length;
    setTheme(themeOrder[nextIndex]);
  }

  const Icon = themeIcons[(theme as keyof typeof themeIcons) ?? "system"] ?? Monitor;

  return (
    <Button variant="ghost" size="icon" onClick={cycleTheme} aria-label="Toggle theme">
      <Icon className="size-4" />
    </Button>
  );
}
