import { useLoadFumenData } from "./fumenData.ts";
import { useState } from "react";
import Dragger from "antd/es/upload/Dragger";
import { Slider, Table, TableProps } from "antd";
import { HighScoreData, useLoadHighScoreData } from "./highScoreData.ts";
import { Rank, RankMap, Style } from "./types.ts";
import { FilterValue } from "antd/es/table/interface";
import { css } from "@emotion/css";
function App() {
  const [fumenFile, setFumenFile] = useState<string | undefined>(undefined);
  const [highScoreFile, setHighScoreFile] = useState<string | undefined>(
    undefined
  );
  const { fumenData } = useLoadFumenData(fumenFile);
  const { highScoreData } = useLoadHighScoreData(highScoreFile, fumenData);
  const [levelRange, setLevelRange] = useState<[number, number]>([0, 19]);
  const [filteredInfo, setFilteredInfo] = useState<
    Record<string, FilterValue | null>
  >({});

  const handleChange: TableProps<HighScoreData>["onChange"] = (
    pagination,
    filters,
    sorter
  ) => {
    console.log("Various parameters", pagination, filters, sorter);
    setFilteredInfo(filters);
  };

  return (
    <>
      <Dragger
        maxCount={1}
        customRequest={(opt) => {
          opt.onSuccess?.({}, undefined);
        }}
        beforeUpload={(f) => {
          f.text().then((t) => {
            setFumenFile(t);
          });
        }}
      >
        upload m_Fumen_List.csv
      </Dragger>
      <Dragger
        maxCount={1}
        customRequest={(opt) => {
          opt.onSuccess?.({}, undefined);
        }}
        beforeUpload={(f) => {
          f.text().then((t) => {
            setHighScoreFile(t);
          });
        }}
      >
        upload t_High_Score.csv
      </Dragger>
      {highScoreData && (
        <Table
          onChange={handleChange}
          size={"small"}
          scroll={{ y: 1000 }}
          rowKey={(r) => r.High_Score_ID}
          pagination={false}
          dataSource={highScoreData}
          columns={[
            {
              title: "Title",
              dataIndex: "Title",
              render: (title: string, record) => {
                const style = record.Style == Style.SP ? "'" : "''";
                return (
                  <span
                    className={css`
                      &[data-difficulty="0"] {
                        color: #00aaaa;
                      }
                      &[data-difficulty="1"] {
                        color: #aaaa00;
                      }
                      &[data-difficulty="2"] {
                        color: #aa0000;
                      }
                      &[data-difficulty="3"] {
                        color: #00aa00;
                      }
                      &[data-difficulty="4"] {
                        color: #aa00aa;
                      }
                    `}
                    data-difficulty={record.Fumen}
                  >
                    {title + style}
                  </span>
                );
              },
            },

            {
              title: "Level",
              dataIndex: "Level",
              filterDropdown: (
                <Slider
                  range
                  min={0}
                  max={19}
                  defaultValue={levelRange}
                  onAfterChange={(value) => setLevelRange(value)}
                />
              ),
              filteredValue: levelRange,
              onFilter: (_, record) => {
                return (
                  Number(record.Level) >= levelRange[0] &&
                  Number(record.Level) <= levelRange[1]
                );
              },
            },
            {
              title: "High Score",
              dataIndex: "Score",
              sorter: (a, b) => a.Score - b.Score,
              render: (score: number, record) => (
                <span
                  className={css`
                    &[data-fc="0"] {
                      color: #00aaaa;
                    }
                    &[data-fc="3"] {
                      color: #00aaaa;
                    }
                    &[data-fc="4"] {
                      color: #00aa00;
                    }
                    &[data-fc="5"] {
                      color: #aaaa00;
                    }
                    &[data-fc="6"] {
                      color: #aaaaaa;
                    }
                  `}
                  data-fc={record.Fumen}
                >
                  {score}
                </span>
              ),
            },
            {
              title: "Rank",
              dataIndex: "Rank",
              filters: RankMap(),
              filteredValue: filteredInfo.Rank,
              onFilter: (value, record) => record.Rank == value,
              render: (rank: Rank) => (
                <>{RankMap().find((v) => v.value == rank)?.text}</>
              ),
            },
          ]}
        />
      )}
    </>
  );
}

export default App;
