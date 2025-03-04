import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

const BackButtonListener = () => {
  const navigate = useNavigate();

  useEffect(() => {
    console.log("Mouse Event");
    const handleMouseBack = (event: any) => {
      if (event.button === 3) {
        navigate(-1);
      } else if (event.button === 4) {
        navigate(1);
      }
    };

    window.addEventListener("mousedown", handleMouseBack);
    return () => {
      window.removeEventListener("mousedown", handleMouseBack);
    };
  }, [navigate]);

  return null;
};

export default BackButtonListener;
