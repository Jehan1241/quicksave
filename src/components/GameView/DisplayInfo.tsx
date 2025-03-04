import { Badge } from "../ui/badge";
import { Button } from "../ui/button";

export function DisplayInfo({ data, tags, companies }: any) {
  return (
    <div className="mt-2 flex h-full flex-col gap-4 overflow-y-auto pr-1 text-left">
      <div className="flex flex-row items-center justify-start">
        <p className="flex flex-col items-start gap-2 text-base font-medium">
          Platform
          <Button className="h-6 rounded-full bg-platformBadge hover:bg-platformBadgeHover text-platformBadgeText">
            {data?.OwnedPlatform}
          </Button>
        </p>
      </div>
      <div className="flex flex-col gap-2 text-base">
        <p className="text-left text-base">Tags</p>
        <div className="flex flex-wrap gap-2 rounded-md text-center">
          {tags.map((items: any, index: any) => (
            <Badge
              className="bg-tagsBadge hover:bg-tagsBadgeHover text-tagsBadgeText"
              key={index}
            >
              {items}
            </Badge>
          ))}
        </div>
      </div>

      <div className="flex flex-col items-start justify-center gap-2 text-base">
        <p>Developers And Publishers</p>
        <div className="flex flex-wrap gap-2 rounded-md text-center">
          {companies.map((items: any, index: any) => (
            <Badge
              className=" bg-devsBadge hover:bg-devsBadgeHover text-devsBadgeText"
              draggable={false}
              key={index}
            >
              {items}
            </Badge>
          ))}
        </div>
      </div>

      <div className="flex flex-col items-start justify-start gap-2">
        Description
        <div className="flex h-full flex-col">
          <p
            dangerouslySetInnerHTML={{ __html: data?.Description }}
            className="text-sm"
          ></p>
        </div>
      </div>
    </div>
  );
}
