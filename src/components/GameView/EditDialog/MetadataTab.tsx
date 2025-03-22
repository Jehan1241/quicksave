import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { TabsContent } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import { CheckedState } from "@radix-ui/react-checkbox";
import { format } from "date-fns";
import { CalendarIcon, Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import MultipleSelector from "@/components/ui/multiple-selector";
import { loadPreferences } from "@/lib/api/GameViewAPI";
import { DateTimePicker } from "@/components/ui/datetime-picker";

export function MetadataTab({ uid, fetchData, tags, companies }: any) {
  const [selectedTags, setSelectedTags] = useState<
    { label: string; value: string }[]
  >(
    Array.isArray(tags)
      ? tags.map((tag: string) => ({ label: tag, value: tag }))
      : []
  );
  const [selectedCompanies, setSelectedCompanies] = useState<
    { label: string; value: string }[]
  >(
    Array.isArray(companies)
      ? companies.map((company: string) => ({ label: company, value: company }))
      : []
  );

  console.log("abawd", companies, selectedCompanies);

  const [customTitleChecked, setCustomTitleChecked] = useState<
    CheckedState | undefined
  >(false);
  const [customTitle, setCustomTitle] = useState("");
  const [customTime, setCustomTime] = useState("");
  const [customTimeOffset, setCustomTimeOffset] = useState("");
  const [customTimeChecked, setCustomTimeChecked] = useState<
    CheckedState | undefined
  >(false);
  const [customTimeOffsetChecked, setCustomTimeOffsetChecked] = useState<
    CheckedState | undefined
  >(false);
  const [customReleaseDate, setCustomReleaseDate] = useState<Date>();
  const [customReleaseDateChecked, setCustomReleaseDateChecked] = useState<
    CheckedState | undefined
  >(false);
  const [customRating, setCustomRating] = useState("");
  const [customRatingChecked, setCustomRatingChecked] = useState<
    CheckedState | undefined
  >(false);
  const [loading, setLoading] = useState(false);
  const [tagOptions, setTagOptions] = useState<
    { label: string; value: string }[]
  >([]);
  const [devOptions, setDevOptions] = useState<
    { label: string; value: string }[]
  >([]);

  const saveClickHandler = () => {
    const selectedTagValues = selectedTags.map(
      (tag: { value: string }) => tag.value
    );
    const selectedDevValues = selectedCompanies.map(
      (dev: { value: string }) => dev.value
    );

    console.log("SELECTED", selectedTagValues);
    const postData = {
      customTitleChecked: customTitleChecked,
      customTitle: customTitle.trim(),
      customTimeChecked: customTimeChecked,
      customTime: customTime.trim() === "" ? "0" : customTime.trim(),
      customTimeOffsetChecked: customTimeOffsetChecked,
      customTimeOffset:
        customTimeOffset.trim() === "" ? "0" : customTimeOffset.trim(),
      customRatingChecked: customRatingChecked,
      customRating: customRating.trim() === "" ? "0" : customRating.trim(),
      customReleaseDateChecked: customReleaseDateChecked,
      customReleaseDate: customReleaseDate
        ? format(customReleaseDate, "yyyy-MM-dd")
        : "", // Empty string if undefined
      UID: uid,
      selectedTags: selectedTagValues,
      selectedDevs: selectedDevValues,
    };
    savePreferences(postData);
  };

  useEffect(() => {
    loadPreferences(
      uid,
      setCustomTime,
      setCustomTimeOffset,
      setCustomRating,
      setCustomTitle,
      setCustomReleaseDate,
      setCustomTitleChecked,
      setCustomTimeChecked,
      setCustomTimeOffsetChecked,
      setCustomReleaseDateChecked,
      setCustomRatingChecked,
      setTagOptions,
      setDevOptions
    );
  }, []);

  const handleCheckboxChange = (
    checked: CheckedState | undefined,
    title: string
  ) => {
    switch (title) {
      case "title":
        setCustomTitleChecked(checked);
        break;
      case "time":
        if (checked == false) {
          setCustomTimeChecked(false);
        } else {
          setCustomTimeChecked(true);
          setCustomTimeOffsetChecked(false);
        }
        break;

      case "timeOffset":
        if (checked == false) {
          setCustomTimeOffsetChecked(false);
        } else {
          setCustomTimeOffsetChecked(true);
          setCustomTimeChecked(false);
        }
        break;
      case "releaseDate":
        setCustomReleaseDateChecked(checked);
        break;
      case "rating":
        setCustomRatingChecked(checked);
        break;
      default:
        break;
    }
  };

  const savePreferences = async (postData: any) => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/SavePreferences`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(postData),
      });
      const json = await response.json();
      if (json.status === "OK") {
        setLoading(false);
        fetchData();
      }
    } catch (error) {
      console.error(error);
      setLoading(false);
    }
  };

  return (
    <TabsContent value="metadata" className="h-full focus:ring-0">
      <div className="flex h-full flex-col justify-between p-2 px-4 focus:outline-none gap-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="flex items-center">
            <Label className="w-60">Use Custom Title</Label>
            <Checkbox
              checked={customTitleChecked}
              onCheckedChange={(checked) =>
                handleCheckboxChange(checked, "title")
              }
              className="mr-10"
            />
            <Input
              disabled={!customTitleChecked}
              id="customTitle"
              value={customTitle}
              onChange={(e) => setCustomTitle(e.target.value)}
            ></Input>
          </div>
          <div className="flex items-center">
            <Label className="w-60">Custom Release Date</Label>
            <Checkbox
              checked={customReleaseDateChecked}
              onCheckedChange={(checked) =>
                handleCheckboxChange(checked, "releaseDate")
              }
              className="mr-10"
            />
            <DateTimePicker
              hideTime={true}
              value={customReleaseDate}
              onChange={setCustomReleaseDate}
              disabled={!customReleaseDateChecked}
            />
          </div>
          <div className="flex items-center">
            <Label className="w-60">Custom Time Played</Label>
            <Checkbox
              checked={customTimeChecked}
              onCheckedChange={(checked) =>
                handleCheckboxChange(checked, "time")
              }
              className="mr-10"
            />
            <Input
              disabled={!customTimeChecked}
              id="customTime"
              value={customTime}
              onChange={(e) => {
                const value = e.target.value;
                // Restrict input to only numbers and up to 2 decimal places
                if (/^\d*\.?\d{0,2}$/.test(value)) {
                  setCustomTime(value);
                }
              }}
            />
          </div>
          <div className="flex items-center">
            <Label className="w-60">Custom Time Offset</Label>
            <Checkbox
              checked={customTimeOffsetChecked}
              onCheckedChange={(checked) =>
                handleCheckboxChange(checked, "timeOffset")
              }
              className="mr-10"
            />
            <Input
              disabled={!customTimeOffsetChecked}
              id="customTimeOffset"
              value={customTimeOffset}
              onChange={(e) => {
                const value = e.target.value;
                if (/^-?\d*\.?\d{0,2}$/.test(value)) {
                  setCustomTimeOffset(value);
                }
              }}
            />
          </div>
          <div className="flex items-center">
            <Label className="w-60">Custom Rating</Label>
            <Checkbox
              checked={customRatingChecked}
              onCheckedChange={(checked) =>
                handleCheckboxChange(checked, "rating")
              }
              className="mr-10"
            />
            <Input
              disabled={!customRatingChecked}
              id="customRating"
              value={customRating}
              onChange={(e) => {
                const value = e.target.value;
                // Restrict input to only numbers and up to 2 decimal places
                if (/^\d*\.?\d{0,2}$/.test(value)) {
                  setCustomRating(value);
                }
              }}
            />
          </div>
          <div className="flex items-center col-span-2 gap-8">
            <Label className="w-20">Tags</Label>
            <MultipleSelector
              options={tagOptions}
              value={selectedTags}
              creatable
              onChange={(selected: any) => {
                setSelectedTags(selected);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Tags"
              emptyIndicator={
                <p className="text-center text-sm">Type to create new tag.</p>
              }
            />
          </div>
          <div className="flex items-center col-span-2 gap-8">
            <Label className="w-20">Developers</Label>
            <MultipleSelector
              options={devOptions}
              value={selectedCompanies}
              creatable
              onChange={(selected: any) => {
                setSelectedCompanies(selected);
              }}
              hidePlaceholderWhenSelected={true}
              placeholder="Select Developers"
              emptyIndicator={
                <p className="text-center text-sm">Type to create new tag.</p>
              }
            />
          </div>
        </div>
        <div className="self-end">
          <Button onClick={saveClickHandler} variant={"dialogSaveButton"}>
            Save {loading && <Loader2 className="animate-spin" />}
          </Button>
        </div>
      </div>
    </TabsContent>
  );
}
