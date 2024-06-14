function DisplayInfo(props) {
  const visible = props.visible;
  const data = props.data;
  const tagsArray = props.tags;
  const companiesArray = props.companies;

  if (visible) {
    return (
      <div className="flex overflow-scroll flex-col p-2 h-2/3 text-left rounded-b-2xl border-2 border-gray-500 backdrop-blur-md bg-black/20">
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">Title</p>
          <p className="p-2 w-2/3">{data.Name}</p>
        </div>
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">Platform</p>
          <p className="p-2 w-2/3">{data.OwnedPlatform}</p>
        </div>
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">
            Time Played
          </p>
          <p className="p-2 w-2/3">{data.TimePlayed} Hours</p>
        </div>
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">Tags</p>
          <p className="p-2 w-2/3">
            {tagsArray.map((items) => (
              <h1>{items}</h1>
            ))}
          </p>
        </div>
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">Tags</p>
          <p className="p-2 w-2/3">
            {tagsArray.map((items) => (
              <h1>{items}</h1>
            ))}
          </p>
        </div>
        <div className="flex flex-row">
          <p className="p-2 mx-2 w-1/2 border-r-2 border-gray-500">
            Devs and Publishers
          </p>
          <p className="p-2 w-2/3">
            {companiesArray.map((items) => (
              <h1>{items}</h1>
            ))}
          </p>
        </div>
        <div className="flex flex-row">
          <p className="p-2">{data.Description}</p>
        </div>
      </div>
    );
  }
}

export default DisplayInfo;
