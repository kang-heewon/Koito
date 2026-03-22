import { type CSSProperties, type ReactNode, useMemo } from "react";
import { motion } from "motion/react";
import useScrollProgress from "../../hooks/useScrollProgress";
import ScrollSection from "./ScrollSection";

export interface RecapSection {
  id: string;
  component: ReactNode;
  gradient: string[];
}

interface RecapLayoutProps {
  sections: RecapSection[];
  title: string;
}

function createGradientStyle(colors: string[]): CSSProperties {
  const [first, second, third, fourth] = [
    colors[0] ?? "rgba(124, 92, 255, 0.28)",
    colors[1] ?? "rgba(83, 193, 255, 0.2)",
    colors[2] ?? "rgba(255, 117, 140, 0.18)",
    colors[3] ?? "rgba(8, 10, 18, 0.98)",
  ];

  return {
    backgroundColor: "var(--color-bg)",
    backgroundImage: [
      `radial-gradient(circle at 18% 20%, ${first} 0%, transparent 42%)`,
      `radial-gradient(circle at 82% 24%, ${second} 0%, transparent 38%)`,
      `radial-gradient(circle at 50% 84%, ${third} 0%, transparent 44%)`,
      `linear-gradient(160deg, ${fourth} 0%, var(--color-bg) 100%)`,
    ].join(", "),
    backgroundPosition: "0% 50%",
    backgroundRepeat: "no-repeat",
    backgroundSize: "180% 180%",
  };
}

export default function RecapLayout({ sections, title }: RecapLayoutProps) {
  const sectionIds = useMemo(() => sections.map((section) => section.id), [sections]);
  const { activeSection, progress, scrollToSection } = useScrollProgress(sectionIds);

  return (
    <div className="relative min-h-screen overflow-x-hidden bg-[var(--color-bg)] text-[var(--color-fg)]">
      <div aria-hidden className="pointer-events-none fixed inset-0 overflow-hidden">
        {sections.map((section, index) => (
          <motion.div
            key={section.id}
            className="absolute inset-0"
            style={createGradientStyle(section.gradient)}
            initial={false}
            animate={{
              opacity: activeSection === index ? 1 : 0,
              scale: activeSection === index ? 1 : 1.04,
              backgroundPosition: ["0% 50%", "100% 50%", "0% 50%"],
            }}
            transition={{
              opacity: { duration: 0.7, ease: "easeInOut" },
              scale: { duration: 0.7, ease: "easeInOut" },
              backgroundPosition: {
                duration: 18,
                ease: "linear",
                repeat: Number.POSITIVE_INFINITY,
              },
            }}
          />
        ))}
        <div className="absolute inset-0 bg-[linear-gradient(180deg,rgba(8,10,18,0.18)_0%,rgba(8,10,18,0.62)_100%)]" />
      </div>

      <div className="pointer-events-none fixed left-0 right-0 top-0 z-20 flex items-center justify-between px-5 pt-5 sm:px-8 sm:pt-8 lg:px-12">
        <div className="pointer-events-auto inline-flex items-center gap-3 rounded-full border border-white/10 bg-[var(--color-bg)]/45 px-4 py-2 text-xs font-semibold uppercase tracking-[0.24em] text-[var(--color-fg)]/78 backdrop-blur-md sm:px-5 sm:py-3">
          <span>{title}</span>
          <span className="h-1 w-16 overflow-hidden rounded-full bg-white/10 sm:w-24">
            <span
              className="block h-full rounded-full bg-[var(--color-primary)] transition-[width] duration-500"
              style={{ width: `${Math.round(progress * 100)}%` }}
            />
          </span>
        </div>
      </div>

      <nav
        aria-label={`${title} sections`}
        className="fixed bottom-5 left-1/2 z-20 flex -translate-x-1/2 items-center gap-3 rounded-full border border-white/10 bg-[var(--color-bg)]/45 px-4 py-3 backdrop-blur-md md:bottom-auto md:left-auto md:right-8 md:top-1/2 md:-translate-y-1/2 md:translate-x-0 md:flex-col md:px-3 md:py-4 lg:right-12"
      >
        {sections.map((section, index) => {
          const isActive = activeSection === index;

          return (
            <button
              key={section.id}
              type="button"
              onClick={() => scrollToSection(index)}
              aria-label={`${index + 1}번째 섹션으로 이동`}
              aria-current={isActive ? "true" : undefined}
              className="group flex h-3 w-3 items-center justify-center rounded-full"
            >
              <span
                className={`block rounded-full transition-all duration-300 ${
                  isActive
                    ? "h-3 w-3 bg-[var(--color-primary)] shadow-[0_0_0_6px_rgba(255,255,255,0.08)]"
                    : "h-2.5 w-2.5 bg-white/45 group-hover:bg-white/72"
                }`}
              />
            </button>
          );
        })}
      </nav>

      <div className="relative z-10">
        {sections.map((section, index) => (
          <ScrollSection
            key={section.id}
            className="px-5 py-20 sm:px-8 sm:py-24 lg:px-12"
            delay={Math.min(index * 0.08, 0.24)}
          >
            <div id={section.id} className="mx-auto w-full max-w-[1100px]">
              {section.component}
            </div>
          </ScrollSection>
        ))}
      </div>
    </div>
  );
}
