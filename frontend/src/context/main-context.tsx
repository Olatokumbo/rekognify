import { createContext, useState, ReactNode, useContext } from "react";

interface MainContextType {
    uuid: string | null;
    setUuid: React.Dispatch<React.SetStateAction<string | null>>;
}

const MainContext = createContext<MainContextType | undefined>(undefined);

export const useMainContext = () => useContext(MainContext)!;

interface MainProviderProps {
    children: ReactNode;
}

const MainProvider = ({ children }: MainProviderProps) => {
    const [uuid, setUuid] = useState<string | null>(null);

    return (
        <MainContext.Provider value={{ uuid, setUuid }}>
            {children}
        </MainContext.Provider>
    );
};

export default MainProvider;