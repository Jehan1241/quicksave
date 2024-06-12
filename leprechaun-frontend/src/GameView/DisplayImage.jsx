import { useState } from "react";
import { GrPrevious, GrNext } from "react-icons/gr";

function DisplayImage(props) {
  const screenshotsArray = props.screenshots;
  const [selectedScreenShot, setSelectedScreenshot] = useState(0);

  const imageNextClickHanlder = () => {
    if (selectedScreenShot < screenshotsArray.length - 1) {
      const newId = selectedScreenShot + 1;
      setSelectedScreenshot(newId);
    }
  };
  const imagePrevClickHandler = () => {
    if (selectedScreenShot > 0) {
      const newId = selectedScreenShot - 1;
      setSelectedScreenshot(newId);
    }
  };

  return (
    <div className="flex flex-col overflow-scroll mt-2 w-[65%] h-2/3 rounded-3xl border-2 border-gray-500 opacity-75 backdrop-blur-3xl">
      <div className="flex justify-center p-4">
        <img
          className="object-scale-down h-[70vh] rounded-3xl"
          src={"/leprechaun-backend/" + screenshotsArray[selectedScreenShot]}
        />
      </div>
      <div className="flex flex-row justify-center mt-1 mb-2">
        <button onClick={imagePrevClickHandler} className="mx-2">
          <GrPrevious size={30} />
        </button>
        <button onClick={imageNextClickHanlder} className="mx-2">
          <GrNext size={30} />
        </button>
      </div>
    </div>
  );
}

export default DisplayImage;
