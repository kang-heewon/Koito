import { average } from "color.js";

export type RecapGradient = [string, string, string, string];

export const wrappedGradientPresets: RecapGradient[] = [
  ["rgba(255, 120, 168, 0.42)", "rgba(129, 91, 255, 0.32)", "rgba(87, 214, 255, 0.26)", "rgba(8, 9, 18, 0.98)"],
  ["rgba(255, 184, 108, 0.38)", "rgba(255, 101, 132, 0.28)", "rgba(147, 83, 255, 0.26)", "rgba(10, 12, 24, 0.98)"],
  ["rgba(118, 255, 182, 0.32)", "rgba(61, 211, 255, 0.28)", "rgba(117, 87, 255, 0.22)", "rgba(8, 12, 24, 0.98)"],
  ["rgba(255, 126, 95, 0.34)", "rgba(255, 73, 171, 0.3)", "rgba(120, 91, 255, 0.24)", "rgba(12, 8, 22, 0.98)"],
  ["rgba(255, 215, 130, 0.34)", "rgba(255, 138, 101, 0.28)", "rgba(255, 94, 98, 0.22)", "rgba(15, 10, 18, 0.98)"],
  ["rgba(137, 247, 254, 0.3)", "rgba(102, 166, 255, 0.24)", "rgba(140, 109, 255, 0.28)", "rgba(8, 11, 20, 0.98)"],
];

export const rewindGradientPresets: RecapGradient[] = [
  ["rgba(93, 211, 255, 0.3)", "rgba(124, 92, 255, 0.26)", "rgba(255, 117, 140, 0.22)", "rgba(6, 10, 20, 0.98)"],
  ["rgba(255, 154, 139, 0.34)", "rgba(255, 109, 170, 0.28)", "rgba(123, 97, 255, 0.22)", "rgba(12, 8, 22, 0.98)"],
  ["rgba(255, 205, 112, 0.32)", "rgba(255, 145, 77, 0.26)", "rgba(255, 94, 98, 0.2)", "rgba(15, 10, 18, 0.98)"],
  ["rgba(118, 255, 182, 0.26)", "rgba(70, 198, 255, 0.24)", "rgba(98, 84, 255, 0.22)", "rgba(8, 14, 20, 0.98)"],
  ["rgba(198, 132, 255, 0.28)", "rgba(123, 104, 238, 0.24)", "rgba(83, 193, 255, 0.2)", "rgba(10, 9, 24, 0.98)"],
  ["rgba(255, 138, 128, 0.32)", "rgba(255, 112, 67, 0.24)", "rgba(255, 215, 0, 0.18)", "rgba(16, 11, 14, 0.98)"],
];

function getPresetByIndex(presets: RecapGradient[], sectionIndex: number) {
  return presets[((sectionIndex % presets.length) + presets.length) % presets.length];
}

export function getWrappedGradient(sectionIndex: number) {
  return [...getPresetByIndex(wrappedGradientPresets, sectionIndex)];
}

export function getRewindGradient(sectionIndex: number) {
  return [...getPresetByIndex(rewindGradientPresets, sectionIndex)];
}

export function mapSectionGradients(
  sectionIds: string[],
  variant: "wrapped" | "rewind",
): Record<string, string[]> {
  const gradientFactory = variant === "wrapped" ? getWrappedGradient : getRewindGradient;

  return Object.fromEntries(
    sectionIds.map((sectionId, sectionIndex) => [sectionId, gradientFactory(sectionIndex)]),
  );
}

export async function getDynamicRewindGradient(
  imageSource: string | null | undefined,
  fallback = rewindGradientPresets[0],
): Promise<string[]> {
  if (!imageSource) {
    return [...fallback];
  }

  try {
    const color = await average(imageSource, { amount: 1 });
    const [red, green, blue] = color as unknown as number[];

    if (
      typeof red !== "number" ||
      typeof green !== "number" ||
      typeof blue !== "number"
    ) {
      return [...fallback];
    }

    return [
      `rgba(${red}, ${green}, ${blue}, 0.44)`,
      `rgba(${Math.max(red - 20, 0)}, ${Math.max(green - 12, 0)}, ${Math.min(blue + 24, 255)}, 0.26)`,
      fallback[2],
      fallback[3],
    ];
  } catch {
    return [...fallback];
  }
}
