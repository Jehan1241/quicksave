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

  const imageTileClickHandler = (index) => {
    setSelectedScreenshot(index);
  };

  return (
    <div className="flex flex-col overflow-y-auto w-[60%] gap-1">
      {/* Main Img Div */}
      <div className="flex justify-center">
        <img
          className="object-scale-down w-full"
          src={
            "http://localhost:8080/screenshots" +
            screenshotsArray[selectedScreenShot]
          }
        />
      </div>
      {/* X-Scrollbar Div */}
      <div className="inline-flex overflow-x-scroll flex-row gap-2 pb-1">
        {screenshotsArray.map((item, index) => (
          <div className={`w-32 min-w-32`}>
            <button onClick={() => imageTileClickHandler(index)}>
              <img
                src={
                  "http://localhost:8080/screenshots" + screenshotsArray[index]
                }
                className={`w-32 ${
                  index === selectedScreenShot ? "border-2" : ""
                }`}
              />
            </button>
          </div>
        ))}
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
