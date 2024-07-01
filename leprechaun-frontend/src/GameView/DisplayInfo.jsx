import { FaSteam } from "react-icons/fa";

function DisplayInfo(props) {
  const data = props.data;
  const tagsArray = props.tags;
  const companiesArray = props.companies;

  return (
    <div className="flex overflow-y-auto flex-col h-2/3 text-left rounded-2xl">
      <div className="flex flex-row justify-start items-center text-xl">
        {data.OwnedPlatform === "Steam" && <FaSteam className="m-1" />}
        <p className="m-1">{data.OwnedPlatform}</p>
      </div>
      <div className="flex flex-col justify-center items-start p-2 text-xl">
        <p>Developers And Publishers</p>
        <p className="m-2 text-base">
          {companiesArray.map((items) => (
            <h1>- {items}</h1>
          ))}
        </p>
      </div>
      <div className="flex flex-col justify-start items-start p-2 text-xl">
        <p className="text-left">Tags</p>
        <div className="text-center">
          {tagsArray.map((items) => (
            <div className="inline-flex justify-center items-center p-2 mx-1 mt-1 h-8 text-sm rounded-md bg-primary">
              {items}
            </div>
          ))}
        </div>
      </div>
      <div className="flex flex-row">
        <p className="p-2">{data.Description}</p>
      </div>
    </div>
  );
}

export default DisplayInfo;
