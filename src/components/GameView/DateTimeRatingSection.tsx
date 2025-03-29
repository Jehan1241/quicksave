import { CalendarDays, Clock, Star } from "lucide-react";

export function DateTimeRatingSection({
  releaseDate,
  rating,
  isWishlist,
  timePlayed,
}: any) {
  return (
    <div className="flex flex-col items-center text-sm xl:ml-auto">
      <div>
        <CalendarDays size={18} className="mb-1 inline" /> {releaseDate}
      </div>
      <div>
        <Star size={18} className="mb-1 inline" /> {rating}
        {isWishlist === 0 && (
          <span>
            <Clock className="mb-1 ml-2 mr-1 inline" size={18} />
            {timePlayed}
          </span>
        )}
      </div>
    </div>
  );
}
