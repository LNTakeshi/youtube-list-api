import { usePapaParse } from "react-papaparse";
import { useEffect, useState } from "react";
import { Difficulty, Rank, Style } from "./types.ts";
import { FumenData } from "./fumenData.ts";

export type HighScoreData = {
  High_Score_ID: string;
  DDR_CODE: string;
  Fumen_ID: string;
  Music_ID: string;
  Title: string;
  Style: Style;
  Fumen: Difficulty;
  Score: number;
  Rank: Rank;
  Full_Combo: number;
  Skill_Point: number;
  Date: Date;
  Level: number;
};

export const useLoadHighScoreData = (
  file?: string,
  fumenData?: FumenData[]
): {
  highScoreData: HighScoreData[];
} => {
  const { readString } = usePapaParse();
  const [data, setData] = useState<HighScoreData[]>([]);
  useEffect(() => {
    if (file) {
      readString<HighScoreData>(file, {
        worker: true,
        header: true,
        complete(results) {
          setData(results.data);
        },
      });
    }
  }, [file, readString]);

  return {
    highScoreData: data.map((d) => {
      return {
        ...d,
        Level: fumenData?.find((f) => f.Fumen_ID === d.Fumen_ID)?.Level ?? 0,
      };
    }),
  };
};
