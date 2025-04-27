import {
  Children,
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";

interface ExePathContextType {
  exePath: string;
}

const ExePathContext = createContext<ExePathContextType | undefined>(undefined);

export const ExePathProvider = ({ children }: { children: ReactNode }) => {
  const [exePath, setExePath] = useState("");

  const getExePath = async () => {
    setExePath(await window.appPaths.exePath());
  };

  useEffect(() => {
    getExePath();
  }, []);

  return (
    <ExePathContext.Provider value={{ exePath }}>
      {children}
    </ExePathContext.Provider>
  );
};

export const useExePathContext = () => {
  const context = useContext(ExePathContext);
  if (!context)
    throw new Error("useExePathContext must be used within a ExePathProvider");
  return context;
};
