import { createContext, useState, ReactNode, useContext } from "react";

interface MainContextType {
    uuid: string | null;
    setUuid: React.Dispatch<React.SetStateAction<string | null>>;
    loading: boolean
    setLoading: React.Dispatch<React.SetStateAction<boolean>>;
}

const MainContext = createContext<MainContextType | undefined>(undefined);

export const useMainContext = () => useContext(MainContext)!;

interface MainProviderProps {
    children: ReactNode;
}

const MainProvider = ({ children }: MainProviderProps) => {
    const [uuid, setUuid] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    return (
        <MainContext.Provider value={{ uuid, setUuid, loading, setLoading }}>
            {children}
        </MainContext.Provider>
    );
};

export default MainProvider;