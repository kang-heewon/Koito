import { ChevronLeft, ChevronRight } from "lucide-react";
import { average } from "color.js";
import { getRewindStats, imageUrl, type RewindStats } from "api/api";
import { useEffect, useState } from "react";
import type { LoaderFunctionArgs } from "react-router";
import { useLoaderData, useLocation, useNavigate } from "react-router";
import Rewind from "~/components/rewind/Rewind";
import { getRewindParams } from "~/utils/utils";

const months = [
  "Full Year",
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

export async function clientLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const fallbackParams = getRewindParams(url.searchParams);
  const parsedYear = Number.parseInt(url.searchParams.get("year") || "", 10);
  const parsedMonth = Number.parseInt(url.searchParams.get("month") || "", 10);
  const year = Number.isNaN(parsedYear) ? fallbackParams.year : parsedYear;
  const month = Number.isNaN(parsedMonth) ? fallbackParams.month : parsedMonth;

  const stats = await getRewindStats({ year, month });
  stats.title = `Your ${month === 0 ? "" : `${months[month]} `}${year} Rewind`;

  return { stats };
}

export function meta({ data }: { data?: { stats: RewindStats } }) {
  const pageTitle = `${data?.stats.title || "Rewind"} - Koito`;

  return [
    { title: pageTitle },
    { name: "description", content: pageTitle },
    { property: "og:title", content: pageTitle },
  ];
}

export default function RewindPage() {
  const { stats } = useLoaderData() as { stats: RewindStats };
  const location = useLocation();
  const navigate = useNavigate();
  const [showTime, setShowTime] = useState(false);
  const [bgColor, setBgColor] = useState("var(--color-bg)");

  const currentParams = new URLSearchParams(location.search);
  const { year, month } = getRewindParams(currentParams);

  useEffect(() => {
    const image = stats.top_artists[0]?.item.image;
    if (!image) {
      return;
    }

    average(imageUrl(image, "small"), { amount: 1 })
      .then((color) => {
        const [red, green, blue] = color as unknown as number[];

        if (
          typeof red !== "number" ||
          typeof green !== "number" ||
          typeof blue !== "number"
        ) {
          return;
        }

        setBgColor(`rgba(${red}, ${green}, ${blue}, 0.4)`);
      })
      .catch(() => {
        setBgColor("var(--color-bg)");
      });
  }, [stats]);

  const updateParams = (params: Record<string, string | null>) => {
    const nextParams = new URLSearchParams(location.search);

    for (const key in params) {
      const value = params[key];

      if (value === null) {
        nextParams.delete(key);
        continue;
      }

      nextParams.set(key, value);
    }

    navigate(`/rewind?${nextParams.toString()}`, { replace: false });
  };

  const navigateMonth = (direction: "prev" | "next") => {
    let nextMonth = month;

    if (direction === "next") {
      nextMonth = month === 12 ? 0 : month + 1;
    } else {
      nextMonth = month === 0 ? 12 : month - 1;
    }

    updateParams({
      year: String(year),
      month: String(nextMonth),
    });
  };

  const navigateYear = (direction: "prev" | "next") => {
    const nextYear = direction === "next" ? year + 1 : year - 1;

    updateParams({
      year: String(nextYear),
      month: String(month),
    });
  };

  return (
    <main
      className="min-h-screen w-full"
      style={{
        background: `linear-gradient(to bottom, ${bgColor}, var(--color-bg) 500px)`,
        transition: "1000ms",
      }}
    >
      <div className="flex flex-col items-start gap-4 sm:items-center">
        <div className="flex w-19/20 flex-col items-start gap-10 px-5 md:px-20 lg:mt-15 lg:flex-row">
          <div className="flex flex-col items-start gap-4">
            <div className="flex flex-col items-start gap-4 py-8">
              <div className="flex items-center justify-around gap-6">
                <button
                  type="button"
                  onClick={() => navigateMonth("prev")}
                  className="cursor-pointer p-2 disabled:text-(--color-fg-tertiary)"
                  disabled={
                    new Date(year, month - 2) > new Date() ||
                    (new Date().getFullYear() === year && month === 1)
                  }
                >
                  <ChevronLeft size={20} />
                </button>
                <p className="w-30 text-center text-xl font-medium">{months[month]}</p>
                <button
                  type="button"
                  onClick={() => navigateMonth("next")}
                  className="cursor-pointer p-2 disabled:text-(--color-fg-tertiary)"
                  disabled={
                    month >= new Date().getMonth() &&
                    year >= new Date().getFullYear()
                  }
                >
                  <ChevronRight size={20} />
                </button>
              </div>

              <div className="flex items-center justify-around gap-6">
                <button
                  type="button"
                  onClick={() => navigateYear("prev")}
                  className="cursor-pointer p-2 disabled:text-(--color-fg-tertiary)"
                  disabled={new Date(year - 1, month) > new Date()}
                >
                  <ChevronLeft size={20} />
                </button>
                <p className="w-30 text-center text-xl font-medium">{year}</p>
                <button
                  type="button"
                  onClick={() => navigateYear("next")}
                  className="cursor-pointer p-2 disabled:text-(--color-fg-tertiary)"
                  disabled={
                    new Date(year + 1, month - 1) > new Date() ||
                    (month === 0 && new Date().getFullYear() === year + 1) ||
                    (new Date().getMonth() === month - 1 &&
                      new Date().getFullYear() === year + 1)
                  }
                >
                  <ChevronRight size={20} />
                </button>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <label htmlFor="show-time-checkbox">Show time listened?</label>
              <input
                id="show-time-checkbox"
                type="checkbox"
                checked={showTime}
                onChange={() => setShowTime((prev) => !prev)}
              />
            </div>
          </div>

          <Rewind stats={stats} includeTime={showTime} />
        </div>
      </div>
    </main>
  );
}
