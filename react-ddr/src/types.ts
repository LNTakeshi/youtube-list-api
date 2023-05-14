export const Style = {
  SP: 0,
  DP: 1,
} as const;

export const StyleMap = () => {
  return [
    {
      text: "SP",
      value: Style.SP,
    },
    {
      text: "DP",
      value: Style.DP,
    },
  ];
};

export type Style = (typeof Style)[keyof typeof Style];

export const Difficulty = {
  Beginner: 0,
  Basic: 1,
  Difficult: 2,
  Expert: 3,
  Challenge: 4,
} as const;

export const DifficultyMap = () => {
  return [
    {
      text: "Beginner",
      value: Difficulty.Beginner,
    },
    {
      text: "Basic",
      value: Difficulty.Basic,
    },
    {
      text: "Difficult",
      value: Difficulty.Difficult,
    },
    {
      text: "Expert",
      value: Difficulty.Expert,
    },
    {
      text: "Challenge",
      value: Difficulty.Challenge,
    },
  ];
};

export type Difficulty = (typeof Difficulty)[keyof typeof Difficulty];

export const Rank = {
  NoPlay: 0,
  Failed: 1,
  D: 2,
  D_Plus: 3,
  C_Minus: 4,
  C: 5,
  C_Plus: 6,
  B_Minus: 7,
  B: 8,
  B_Plus: 9,
  A_Minus: 10,
  A: 11,
  A_Plus: 12,
  AA_Minus: 13,
  AA: 14,
  AA_Plus: 15,
  AAA: 16,
} as const;

export const RankMap = () => {
  return [
    {
      text: "No Play",
      value: Rank.NoPlay,
    },
    {
      text: "Failed",
      value: Rank.Failed,
    },
    {
      text: "D",
      value: Rank.D,
    },
    {
      text: "D+",
      value: Rank.D_Plus,
    },
    {
      text: "C-",
      value: Rank.C_Minus,
    },
    {
      text: "C",
      value: Rank.C,
    },
    {
      text: "C+",
      value: Rank.C_Plus,
    },
    {
      text: "B-",
      value: Rank.B_Minus,
    },
    {
      text: "B",
      value: Rank.B,
    },
    {
      text: "B+",
      value: Rank.B_Plus,
    },
    {
      text: "A-",
      value: Rank.A_Minus,
    },
    {
      text: "A",
      value: Rank.A,
    },
    {
      text: "A+",
      value: Rank.A_Plus,
    },
    {
      text: "AA-",
      value: Rank.AA_Minus,
    },
    {
      text: "AA",
      value: Rank.AA,
    },
    {
      text: "AA+",
      value: Rank.AA_Plus,
    },
    {
      text: "AAA",
      value: Rank.AAA,
    },
  ];
};

export type Rank = (typeof Rank)[keyof typeof Rank];

export const FullCombo = {
  None: 0,
  FC: 3,
  GFC: 4,
  PFC: 5,
  MFC: 6,
} as const;

export const FullComboMap = () => {
  return [
    {
      text: "None",
      value: FullCombo.None,
    },
    {
      text: "FC",
      value: FullCombo.FC,
    },
    {
      text: "GFC",
      value: FullCombo.GFC,
    },
    {
      text: "PFC",
      value: FullCombo.PFC,
    },
    {
      text: "MFC",
      value: FullCombo.MFC,
    },
  ];
};

export type FullCombo = (typeof FullCombo)[keyof typeof FullCombo];
