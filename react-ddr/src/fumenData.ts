import { usePapaParse } from "react-papaparse";
import { useEffect, useState } from "react";

export type FumenData = {
  Fumen_ID: string;
  Music_ID: string;
  Initial: number;
  Title: string;
  Version: number;
  Deleted: number;
  Style: number;
  Fumen: number;
  Level: number;
  Shock: number;
  BPM: number;
  Stream: number;
  Voltage: number;
  Air: number;
  Freeze: number;
  Chaos: number;
};

export const useLoadFumenData = (file?: string) => {
  const { readString } = usePapaParse();
  const [data, setData] = useState<FumenData[]>([]);

  useEffect(() => {
    if (file) {
      readString<FumenData>(file, {
        worker: true,
        header: true,
        complete(results) {
          setData(results.data);
        },
      });
    }
  }, [file, readString]);

  return {
    fumenData: data,
  };
};
