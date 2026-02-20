import TopItemList from "~/components/TopItemList";
import ChartLayout from "./ChartLayout";
import { useLoaderData, type LoaderFunctionArgs } from "react-router";
import { type Artist as ArtistItem, type PaginatedResponse } from "api/api";

export async function clientLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const page = url.searchParams.get("page") || "1";
  url.searchParams.set('page', page)

  const res = await fetch(
    `/apis/web/v1/top-artists?${url.searchParams.toString()}`
  );
  if (!res.ok) {
    throw new Response("Failed to load top artists", { status: 500 });
  }

  const top_artists: PaginatedResponse<ArtistItem> = await res.json();
  return { top_artists };
}

export default function ArtistChart() {
  const { top_artists: initialData } = useLoaderData<{ top_artists: PaginatedResponse<ArtistItem> }>();

  return (
    <ChartLayout
      title="Top Artists"
      initialData={initialData}
      endpoint="chart/top-artists"
      render={({ data, page, onNext, onPrev }) => (
        <div className="flex flex-col gap-5">
        <div className="flex gap-15 mx-auto">
          <button type="button" className="default" onClick={onPrev} disabled={page <= 1}>
            Prev
          </button>
          <button type="button" className="default" onClick={onNext} disabled={!data.has_next_page}>
            Next
          </button>
        </div>
          <TopItemList
            separators
            data={data}
            className="w-[400px] sm:w-[600px]"
            type="artist"
          />
          <div className="flex gap-15 mx-auto">
            <button type="button" className="default" onClick={onPrev} disabled={page <= 1}>
              Prev
            </button>
            <button type="button" className="default" onClick={onNext} disabled={!data.has_next_page}>
              Next
            </button>
          </div>
        </div>
      )}
    />
  );
}
