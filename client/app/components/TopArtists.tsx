import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { getTopArtists, type GetItemsArgs } from "api/api";
import { Link } from "react-router";
import TopListSkeleton from "./skeletons/TopListSkeleton";
import TopItemList from "./TopItemList";

interface Props {
  limit: number;
  period: string;
  artistId?: number;
  albumId?: number;
}

export default function TopArtists(props: Props) {
  const { isPending, isError, data, error, isFetching } = useQuery({
    queryKey: [
      "top-artists",
      { limit: props.limit, period: props.period, page: 0 },
    ],
    queryFn: ({ queryKey }) => getTopArtists(queryKey[1] as GetItemsArgs),
    placeholderData: keepPreviousData,
  });
  if (isPending || isFetching) {
    return (
      <div className="w-[300px]">
        <h2>Top Artists</h2>
        <TopListSkeleton numItems={props.limit} />
      </div>
    );
  } else if (isError) {
    return (
      <div className="w-[300px]">
        <h2>Top Artists</h2>
        <p className="error">Error: {error.message}</p>
      </div>
    );
  }

  if (!data?.items) return null;

  return (
    <div className="w-[300px]">
      <h2 className="hover:underline">
        <Link to={`/chart/top-artists?period=${props.period}`}>
          Top Artists
        </Link>
      </h2>
      <div className="max-w-[300px]">
        <TopItemList type="artist" data={data} />
        {data.items.length < 1 ? "No artists found for this period." : ""}
      </div>
    </div>
  );
}
