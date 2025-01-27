import React, { useState, useEffect } from "react";
import { useSortContext } from "@/SortContext";
import { Dialog, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { DialogContent, DialogDescription } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import { Loader2 } from "lucide-react";

export default function Integrations() {
    const { isIntegrationsDialogOpen, setIsIntegrationsDialogOpen } = useSortContext();
    const [steamIDEmpty, setSteamIDEmpty] = useState(false);
    const [apiKeyEmpty, setAPIKeyEmpty] = useState(false);
    const [loading, setLoading] = useState(false);
    const [steamID, setSteamID] = useState("");
    const [apiKey, setApiKey] = useState("");
    const [steamError, setSteamError] = useState(null);
    const { toast } = useToast();
    const [npsso, setNpsso] = useState("");
    const [npssoEmpty, setNpssoEmpty] = useState(false);
    const [psnLoading, setPsnLoading] = useState(false);
    const [psnError, setPsnError] = useState(null);
    const [psnGamesNotMatched, setPsnGamesNotMatched] = useState<string[]>([]);

    useEffect(() => {
        if (loading === true) {
            toast({
                variant: "default",
                title: "Steam Integration Started!",
                description: "You can safely leave this page now.",
            });
        }
        if (steamError === true) {
            toast({
                variant: "destructive",
                title: "Steam Integration Error!",
                description: "Please check your credentials and try again.",
            });
            setSteamError(null);
        } else if (steamError === false) {
            toast({
                variant: "default",
                title: "Library Integrated!",
                description: "Your steam library has been successfully integrated.",
            });
            setSteamError(null);
        }
        if (psnLoading === true) {
            toast({
                variant: "default",
                title: "PSN Integration Started!",
                description: "You can safely leave this page now.",
            });
        }
        if (psnError === true) {
            toast({
                variant: "destructive",
                title: "PSN Integration Error!",
                description: "Please check your npsso and try again.",
            });
            setPsnError(null);
        } else if (psnError === false) {
            toast({
                variant: "default",
                title: "Library Integrated!",
                description: "Your PSN library has been successfully integrated.",
            });
            setPsnError(null);
        }
    }, [steamError, loading, psnLoading, psnError]);

    const importSteamLibrary = async () => {
        const SteamIdElement = document.getElementById("steamid") as HTMLInputElement;
        const APIKeyElement = document.getElementById("apikey") as HTMLInputElement;

        const SteamID = SteamIdElement.value;
        const APIkey = APIKeyElement.value;

        if (!SteamID) {
            setSteamIDEmpty(true);
        }
        if (!APIkey) {
            setAPIKeyEmpty(true);
        }

        if (SteamID && APIkey) {
            setSteamIDEmpty(false);
            setAPIKeyEmpty(false);
            console.log("Sending Import Steam Library");
            setLoading(true);
            try {
                const response = await fetch("http://localhost:8080/SteamImport", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({ SteamID: SteamID.trim(), APIkey: APIkey.trim() }),
                });
                const resp = await response.json();
                console.log("Steam Errors ? :", resp.error);
                setSteamError(resp.error);
            } catch (error) {
                console.error("Error:", error);
                setLoading(false);
            }
            setLoading(false);
        }
    };

    const getNpsso = async () => {
        console.log("Getting Npsso");
        try {
            const response = await fetch(`http://localhost:8080/Npsso`);
            const json = await response.json();
            console.log(json.Npsso);
            setNpsso(json.Npsso);
        } catch (error) {
            console.error(error);
        }
    };

    const getSteamCreds = async () => {
        console.log("Getting Steam Creds");
        try {
            const response = await fetch(`http://localhost:8080/SteamCreds`);
            const json = await response.json();
            setSteamID(json.SteamCreds[0]);
            setApiKey(json.SteamCreds[1]);
        } catch (error) {
            console.error(error);
        }
    };

    const importPlaystationLibrary = async () => {
        const npssoElement = document.getElementById("npsso") as HTMLInputElement;
        const npsso = npssoElement.value;
        if (!npsso) {
            setNpssoEmpty(true);
            return;
        } else {
            console.log("Sending PlayStation Import Req", npsso);
            setPsnLoading(true);
            try {
                const response = await fetch("http://localhost:8080/PlayStationImport", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        npsso: npsso,
                        clientID: "bg50w140115zmfq2pi0uc0wujj9pn6",
                        clientSecret: "1nk95mh97tui5t1ct1q5i7sqyfmqvd",
                    }),
                });
                const resp = await response.json();
                console.log("PSN Error : ", resp.error);
                console.log("PSN Games Not Matched : ", resp.gamesNotMatched);
                setPsnError(resp.error);
                if (resp.error === false) {
                    setPsnGamesNotMatched(resp.gamesNotMatched);
                }
            } catch (error) {
                console.error("Error:", error);
                setPsnLoading(false);
            }
            setPsnLoading(false);
        }
    };

    useEffect(() => {
        getSteamCreds();
        getNpsso();
    }, []);

    return (
        <Dialog open={isIntegrationsDialogOpen} onOpenChange={setIsIntegrationsDialogOpen}>
            <DialogContent
                className="flex h-[400px] w-[500px] max-w-[500px] select-none flex-col lg:h-[500px] lg:w-[900px] lg:max-w-[900px]"
                tabIndex={-1}
            >
                <DialogHeader>
                    <DialogTitle>Configure Integrations</DialogTitle>
                    <DialogDescription>
                        You can configure external libraries here. This will automatically import
                        your libraries.
                    </DialogDescription>
                </DialogHeader>

                <div className="h-full w-full text-sm">
                    <Tabs defaultValue="steam" className="flex h-full w-full flex-col">
                        <TabsList className="h-10 w-fit max-w-none">
                            <TabsTrigger value="steam">Steam</TabsTrigger>
                            <TabsTrigger value="playstation">PlayStation</TabsTrigger>
                        </TabsList>
                        <div className="relative flex-1">
                            <TabsContent
                                value="steam"
                                className="absolute inset-0 flex flex-col border-t-0 p-2"
                                tabIndex={-1}
                            >
                                <div className="flex-1">
                                    <div className="flex w-full flex-col items-center justify-center gap-2 text-sm">
                                        <div className="flex w-full items-center">
                                            <label
                                                className={`w-40 ${steamIDEmpty ? "text-destructive" : null}`}
                                            >
                                                {steamIDEmpty && "*"}Steam ID
                                            </label>
                                            <Input
                                                value={steamID}
                                                onChange={(e) => setSteamID(e.target.value)}
                                                id="steamid"
                                                className="h-8 w-full"
                                            />
                                        </div>
                                        <div className="flex w-full items-center">
                                            <label
                                                className={`w-40 ${apiKeyEmpty ? "text-destructive" : null}`}
                                            >
                                                {apiKeyEmpty && "*"}Steam API Key
                                            </label>
                                            <Input
                                                id="apikey"
                                                className="h-8 w-full"
                                                value={apiKey}
                                                onChange={(e) => setApiKey(e.target.value)}
                                            />
                                        </div>
                                    </div>
                                </div>
                                <div className="flex justify-end">
                                    <Button
                                        variant="secondary"
                                        className="w-60"
                                        onClick={importSteamLibrary}
                                        disabled={loading}
                                    >
                                        Import Steam Library{" "}
                                        {loading && <Loader2 className="animate-spin" />}
                                    </Button>
                                </div>
                            </TabsContent>

                            <TabsContent
                                tabIndex={-1}
                                value="playstation"
                                className="absolute inset-0 p-2"
                            >
                                <div className="flex h-full flex-col gap-2">
                                    <div className="flex w-full items-center">
                                        <label
                                            className={`w-40 ${npssoEmpty ? "text-destructive" : null}`}
                                        >
                                            {npssoEmpty && "*"}Npsso
                                        </label>
                                        <Input
                                            value={npsso}
                                            onChange={(e) => setNpsso(e.target.value)}
                                            id="npsso"
                                            className="h-8 w-full"
                                        />
                                    </div>
                                    <div className="flex h-full flex-col overflow-y-auto rounded-md border border-border p-2 text-sm">
                                        <div className="flex h-full flex-col">
                                            {psnGamesNotMatched.length > 0 ? (
                                                psnGamesNotMatched.map((game, index) => (
                                                    <div key={index}>
                                                        <p>{game}</p>
                                                    </div>
                                                ))
                                            ) : (
                                                <div className="text-center text-xs">
                                                    The Playstation API is not perfect and sometimes
                                                    games fail to be added. Games that could not be
                                                    added will be listed here.
                                                </div>
                                            )}
                                        </div>
                                    </div>

                                    <div className="flex justify-end">
                                        <Button
                                            variant="secondary"
                                            className="w-60"
                                            onClick={importPlaystationLibrary}
                                            disabled={psnLoading}
                                        >
                                            Import PlayStation Library
                                            {psnLoading && <Loader2 className="animate-spin" />}
                                        </Button>
                                    </div>
                                </div>
                            </TabsContent>
                        </div>
                    </Tabs>
                </div>
            </DialogContent>
        </Dialog>
    );
}
