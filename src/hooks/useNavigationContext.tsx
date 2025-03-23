import React, {
  createContext,
  useState,
  useContext,
  ReactNode,
  useEffect,
} from "react";
import { useLocation } from "react-router-dom";

interface NavigationContextType {
  lastLibraryPath: string;
}

const NavigationContext = createContext<NavigationContextType | undefined>(
  undefined
);

export const NavigationProvider = ({ children }: { children: ReactNode }) => {
  const location = useLocation();
  const [lastLibraryPath, setLastLibraryPath] = useState<string>("/library");

  useEffect(() => {
    if (
      location.pathname.startsWith("/library") ||
      location.pathname.startsWith("/wishlist") ||
      location.pathname.startsWith("/installed")
    ) {
      setLastLibraryPath(location.pathname);
    }
  }, [location.pathname]);

  return (
    <NavigationContext.Provider value={{ lastLibraryPath }}>
      {children}
    </NavigationContext.Provider>
  );
};

// Custom hook to use the navigation context
export const useNavigationContext = (): NavigationContextType => {
  const context = useContext(NavigationContext);
  if (!context)
    throw new Error(
      "useNavigationContext must be used within a NavigationProvider"
    );
  return context;
};
