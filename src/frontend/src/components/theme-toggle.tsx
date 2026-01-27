"use client";

import { Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";

const themeOrder = ["system", "light", "dark"] as const;

const themeIcons = {
  system: Monitor,
  light: Sun,
  dark: Moon,
} as const;

export function ThemeToggle() {
  const { theme, resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  function cycleTheme() {
    const currentIndex = themeOrder.indexOf(theme as (typeof themeOrder)[number]);
    let nextIndex = (currentIndex === -1 ? 0 : currentIndex + 1) % themeOrder.length;
    // Skip if the next theme resolves to the same appearance (e.g. "system" → "light" when OS is light)
    if (resolvedTheme && themeOrder[nextIndex] === resolvedTheme) {
      nextIndex = (nextIndex + 1) % themeOrder.length;
    }
    setTheme(themeOrder[nextIndex]);
  }

  const Icon = themeIcons[(theme as keyof typeof themeIcons) ?? "system"] ?? Monitor;

  return (
    <Button variant="ghost" size="icon" onClick={cycleTheme} aria-label="Toggle theme">
      {mounted ? <Icon className="size-4" /> : <span className="size-4" />}
    </Button>
  );
}
