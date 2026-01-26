import { describe, expect, it } from "vitest";
import {
  calculateDurationMonths,
  calculateTotalTenure,
  formatDate,
  formatDuration,
  parseToDate,
} from "./utils";

describe("formatDate", () => {
  it("formats ISO dates to human readable", () => {
    expect(formatDate("2020-01-15")).toBe("Jan 2020");
    expect(formatDate("2023-12-01")).toBe("Dec 2023");
  });

  it("returns null for null/undefined input", () => {
    expect(formatDate(null)).toBeNull();
    expect(formatDate(undefined)).toBeNull();
  });

  it("handles Present/current", () => {
    expect(formatDate("Present")).toBe("Present");
    expect(formatDate("present")).toBe("Present");
    expect(formatDate("current")).toBe("Present");
  });

  it("returns human-readable strings as-is", () => {
    expect(formatDate("Jan 2020")).toBe("Jan 2020");
  });

  it("returns null for invalid data", () => {
    expect(formatDate('{"invalid": "json"}')).toBeNull();
    expect(formatDate("-0001-11-30")).toBeNull();
  });
});

describe("parseToDate", () => {
  it("parses ISO format dates", () => {
    const date = parseToDate("2020-01-15");
    expect(date).not.toBeNull();
    expect(date?.getFullYear()).toBe(2020);
    expect(date?.getMonth()).toBe(0); // January
  });

  it("parses Month Year format", () => {
    const date = parseToDate("Jan 2020");
    expect(date).not.toBeNull();
    expect(date?.getFullYear()).toBe(2020);
    expect(date?.getMonth()).toBe(0);

    const dec = parseToDate("December 2023");
    expect(dec?.getFullYear()).toBe(2023);
    expect(dec?.getMonth()).toBe(11);
  });

  it("parses year-only format", () => {
    const date = parseToDate("2020");
    expect(date).not.toBeNull();
    expect(date?.getFullYear()).toBe(2020);
  });

  it("returns current date for Present", () => {
    const date = parseToDate("Present");
    expect(date).not.toBeNull();
    const now = new Date();
    expect(date?.getFullYear()).toBe(now.getFullYear());
  });

  it("returns null for invalid data", () => {
    expect(parseToDate(null)).toBeNull();
    expect(parseToDate(undefined)).toBeNull();
    expect(parseToDate('{"invalid"}')).toBeNull();
    expect(parseToDate("-0001-11-30")).toBeNull();
  });
});

describe("calculateDurationMonths", () => {
  it("calculates months between two dates", () => {
    expect(calculateDurationMonths("2020-01-01", "2020-06-01", false)).toBe(5);
    expect(calculateDurationMonths("2020-01-01", "2021-01-01", false)).toBe(12);
    expect(calculateDurationMonths("2020-01-01", "2022-07-01", false)).toBe(30);
  });

  it("handles human-readable dates", () => {
    expect(calculateDurationMonths("Jan 2020", "Jun 2020", false)).toBe(5);
    expect(calculateDurationMonths("Jan 2020", "Jan 2021", false)).toBe(12);
  });

  it("returns null when dates cannot be parsed", () => {
    expect(calculateDurationMonths(null, "2020-06-01", false)).toBeNull();
    expect(calculateDurationMonths("2020-01-01", null, false)).toBeNull();
  });

  it("uses current date for current positions", () => {
    const result = calculateDurationMonths("2020-01-01", null, true);
    expect(result).not.toBeNull();
    expect(result).toBeGreaterThan(0);
  });
});

describe("formatDuration", () => {
  it("formats months as years and months", () => {
    expect(formatDuration(1)).toBe("1 mo");
    expect(formatDuration(11)).toBe("11 mos");
    expect(formatDuration(12)).toBe("1 yr");
    expect(formatDuration(13)).toBe("1 yr 1 mo");
    expect(formatDuration(24)).toBe("2 yrs");
    expect(formatDuration(25)).toBe("2 yrs 1 mo");
    expect(formatDuration(30)).toBe("2 yrs 6 mos");
  });

  it("handles zero/negative months", () => {
    expect(formatDuration(0)).toBe("< 1 mo");
    expect(formatDuration(-1)).toBe("< 1 mo");
  });
});

describe("calculateTotalTenure", () => {
  it("calculates total for single experience", () => {
    const result = calculateTotalTenure([
      { startDate: "2020-01-01", endDate: "2020-06-01", isCurrent: false },
    ]);
    expect(result).toBe(5);
  });

  it("sums non-overlapping experiences", () => {
    const result = calculateTotalTenure([
      { startDate: "2020-01-01", endDate: "2020-06-01", isCurrent: false },
      { startDate: "2021-01-01", endDate: "2021-06-01", isCurrent: false },
    ]);
    expect(result).toBe(10); // 5 + 5 months
  });

  it("merges overlapping experiences", () => {
    // Two overlapping ranges: Jan-Jun 2020 and Mar-Aug 2020 should merge to Jan-Aug 2020 (7 months)
    const result = calculateTotalTenure([
      { startDate: "2020-01-01", endDate: "2020-06-01", isCurrent: false },
      { startDate: "2020-03-01", endDate: "2020-08-01", isCurrent: false },
    ]);
    expect(result).toBe(7);
  });

  it("handles adjacent experiences", () => {
    // Jun-Dec 2020 and Dec 2020-Jun 2021 are adjacent/continuous
    const result = calculateTotalTenure([
      { startDate: "2020-06-01", endDate: "2020-12-01", isCurrent: false },
      { startDate: "2020-12-01", endDate: "2021-06-01", isCurrent: false },
    ]);
    expect(result).toBe(12);
  });

  it("returns null when no valid experiences", () => {
    expect(calculateTotalTenure([])).toBeNull();
    expect(calculateTotalTenure([{ startDate: null, endDate: null, isCurrent: false }])).toBeNull();
  });

  it("handles current positions", () => {
    const result = calculateTotalTenure([
      { startDate: "2020-01-01", endDate: null, isCurrent: true },
    ]);
    expect(result).not.toBeNull();
    expect(result).toBeGreaterThan(0);
  });
});
