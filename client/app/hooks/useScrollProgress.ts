import { useCallback, useEffect, useState } from "react";

interface UseScrollProgressResult {
  activeSection: number;
  progress: number;
  scrollToSection: (sectionIndex: number) => void;
}

function clampProgress(value: number) {
  if (Number.isNaN(value) || value < 0) {
    return 0;
  }

  if (value > 1) {
    return 1;
  }

  return value;
}

export function useScrollProgress(sectionIds: string[]): UseScrollProgressResult {
  const [activeSection, setActiveSection] = useState(0);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    if (sectionIds.length === 0) {
      setActiveSection(0);
      setProgress(0);
      return;
    }

    const updateProgress = () => {
      const maxScroll = document.documentElement.scrollHeight - window.innerHeight;

      if (maxScroll <= 0) {
        setProgress(0);
        return;
      }

      setProgress(clampProgress(window.scrollY / maxScroll));
    };

    const elements = sectionIds
      .map((sectionId) => {
        const sectionAnchor = document.getElementById(sectionId);
        const sectionElement = sectionAnchor?.closest("section");

        if (!(sectionElement instanceof HTMLElement)) {
          return null;
        }

        return {
          sectionId,
          element: sectionElement,
        };
      })
      .filter(
        (
          sectionEntry,
        ): sectionEntry is {
          sectionId: string;
          element: HTMLElement;
        } => sectionEntry !== null,
      );

    const elementIdMap = new Map(
      elements.map((sectionEntry) => [sectionEntry.element, sectionEntry.sectionId]),
    );

    const observer = new IntersectionObserver(
      (entries) => {
        const visibleEntries = entries
          .filter((entry) => entry.isIntersecting)
          .sort((entryA, entryB) => entryB.intersectionRatio - entryA.intersectionRatio);

        const nextEntry = visibleEntries[0];

        if (!nextEntry) {
          return;
        }

        const nextSectionId = elementIdMap.get(nextEntry.target as HTMLElement);
        const nextIndex = sectionIds.findIndex((sectionId) => sectionId === nextSectionId);

        if (nextIndex >= 0) {
          setActiveSection(nextIndex);
        }
      },
      {
        root: null,
        threshold: [0.2, 0.4, 0.6, 0.8],
        rootMargin: "-15% 0px -25% 0px",
      },
    );

    elements.forEach((sectionEntry) => {
      observer.observe(sectionEntry.element);
    });

    updateProgress();
    window.addEventListener("scroll", updateProgress, { passive: true });
    window.addEventListener("resize", updateProgress);

    return () => {
      observer.disconnect();
      window.removeEventListener("scroll", updateProgress);
      window.removeEventListener("resize", updateProgress);
    };
  }, [sectionIds]);

  const scrollToSection = useCallback(
    (sectionIndex: number) => {
      const targetId = sectionIds[sectionIndex];

      if (!targetId) {
        return;
      }

      document.getElementById(targetId)?.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    },
    [sectionIds],
  );

  return { activeSection, progress, scrollToSection };
}

export default useScrollProgress;
