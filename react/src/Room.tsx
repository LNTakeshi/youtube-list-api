import { blue, grey } from '@ant-design/colors';
import { css } from '@linaria/core';
import {
  Button,
  Col,
  Form,
  Grid,
  Input,
  Row,
  Space,
  Table,
  TimePicker,
  Typography,
  notification,
  Checkbox
} from 'antd';
import { Content } from 'antd/lib/layout/layout';
import Column from 'antd/lib/table/Column';
import AntdLink from 'antd/lib/typography/Link';
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import moment from 'moment';
import React, {
  ReactElement,
  useContext,
  useEffect,
  useRef,
  useState
} from 'react';
import { useCookies } from 'react-cookie';
import { SiNiconico, SiSpotify, SiTwitter, SiYoutube } from 'react-icons/si';
import { Link, useLocation, useMatch, useNavigate } from 'react-router-dom';
import { ColorContext } from './App';
import { ChangeColorButton, useChangeColor } from './Home';

const { Text } = Typography;

const contentCss = css`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const headerCss = css`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const paddingCss = css`
  padding: 8px;
  padding-bottom: 0px;
`;

const deletedRow = css`
  background-color: ${grey[5]};
  color: ${grey[0]};
`;

const currentRowLight = css`
  background-color: ${blue[1]};
`;

const currentRowDark = css`
  background-color: ${blue[9]};
`;

const right = css`
  display: flex !important;
  justify-content: right;
  text-align: right;
`;

const paddingBottom = css`
  padding-bottom: 0px;
`;

type Data = {
  key: string;
  time: string;
  url: string;
  title: string;
  username: string;
  length: string;
  deleted: boolean;
  removable: boolean;
  icon: ReactElement;
};
type GetResponse = {
  data: Data[];
  privateInfo: {
    masterId: string;
    uuid: string;
  };
  info: {
    currentIndex: number;
    needUpdate: boolean;
    calcTime: string;
  };
};

const ROOT = () => {
  if (window.location.href.includes('3000')) {
    return 'http://localhost:8080';
  }
  return '';
};

let lastUpdateDate = moment(0);
let lastSendTime = moment();
let d: GetResponse | undefined = undefined;

const Room = () => {
  const [sending, setSending] = useState(false);
  const [stop, setStop] = useState(false);
  const [data, setData] = useState<GetResponse | undefined>(undefined);
  const match = useMatch(process.env.PUBLIC_URL + '/room/:roomId');
  const navigate = useNavigate();
  if (match == null) {
    navigate(process.env.PUBLIC_URL + '/');
  }
  const location = useLocation();
  const forceReload = location.search == '?u';
  const masterIdStr = `master-id-${match!.params.roomId}`;
  const [cookies, setCookie] = useCookies(['name', 'uuid', masterIdStr]);
  const [form] = Form.useForm();
  const changeColor = useChangeColor();
  const masterId = cookies[masterIdStr] as string | undefined;

  const send = (v: any) => {
    lastSendTime = moment();
    setSending(true);
    setStop(false);
    localStorage.setItem('name', v.name);
    let start = moment(0).utc();
    let end = moment(0).utc();
    if (v.start) {
      start = v.start[0];
      end = v.start[1];
    }
    const params = new URLSearchParams();
    console.log(v);
    params.append('username', v.name ?? '');
    params.append('start', start.format('HH:mm:ss'));
    params.append('end', end.format('HH:mm:ss'));
    params.append(
      'title',
      (v.hidden == true ? '[HIDDEN]' : '') + (v.title ?? '')
    );
    params.append('url', v.url);
    params.append('uuid', cookies.uuid);
    params.append('room_id', match!.params.roomId!);
    axios
      .post(`${ROOT()}/youtube-list/api/youtubelist/send`, params)
      .then(() => {
        getList();
        form.resetFields(['url', 'title', 'start', 'hidden']);
      })
      .catch((e: AxiosError) => {
        notification['error']({
          message: 'Error',
          description: e.response?.data.error ?? e.message
        });
        setSending(false);
      })
      .finally(() => {
        setSending(false);
      });
  };
  const getList = () => {
    if (
      !forceReload &&
      moment().diff(lastSendTime, 'minute') >= 60 &&
      masterId == undefined
    ) {
      setStop(true);
      return;
    }

    const requestOpt: AxiosRequestConfig = {
      url: ROOT() + '/youtube-list/api/youtubelist/getList',
      method: 'GET',
      params: {
        room_id: match!.params.roomId,
        uuid: cookies.uuid,
        master_id: cookies[masterIdStr],
        lastUpdateDate: lastUpdateDate.format()
      }
    };
    axios(requestOpt).then((res: AxiosResponse<GetResponse>) => {
      lastUpdateDate = moment();
      res.data.data.forEach((v) => {
        v.key = v.time + v.url;
      });
      if (res.data.privateInfo.uuid) {
        setCookie('uuid', res.data.privateInfo.uuid, {
          maxAge: 30 * 24 * 60 * 60
        });
      }
      if (res.data.privateInfo.masterId) {
        setCookie(masterIdStr, res.data.privateInfo.masterId, {
          maxAge: 30 * 24 * 60 * 60
        });
      }
      if (res.data.info.needUpdate) {
        res.data.data = res.data.data.reverse();
      } else if (d) {
        res.data.data = d.data;
      }

      const difftime = res.data.data
        .slice(0, res.data.data.length - res.data.info.currentIndex - 1)
        .reduce(
          (acc, v) => {
            if (v.deleted) {
              return acc;
            }
            const times = v.length.split(':');
            acc[0] += parseInt(times[0], 10);
            acc[1] += parseInt(times[1], 10);
            acc[2] += parseInt(times[2], 10);
            return acc;
          },
          [0, 0, 0]
        );
      const calctime = moment()
        .add(difftime[0], 'hours')
        .add(difftime[1], 'minutes')
        .add(difftime[2], 'seconds');
      res.data.info.calcTime = `約${calctime.diff(
        moment(),
        'minutes'
      )}分後${calctime.format('(DD日HH:mm頃)')}`;

      res.data.data.forEach((e) => {
        if (e.url.includes('twitter')) {
          e.icon = <SiTwitter />;
        } else if (e.url.includes('youtube')) {
          e.icon = <SiYoutube />;
        } else if (e.url.includes('nicovideo')) {
          e.icon = <SiNiconico />;
        } else if (e.url.includes('spotify')) {
          e.icon = <SiSpotify />;
        }
      });

      d = res.data;
      setData(res.data);
    });
  };

  useEffect(() => {
    const interval = () => {
      getList();
    };
    const t = setInterval(interval, 30000);
    getList();
    return () => clearInterval(t);
  }, []);

  return (
    <Content className={contentCss}>
      <div className={headerCss}>
        <Row gutter={16} className={paddingBottom} align="middle">
          <Col>
            <ChangeColorButton onClick={changeColor} />
          </Col>
          <Col>
            <Link to={process.env.PUBLIC_URL + '/'}>TOPに戻る</Link>
          </Col>
          <Col>
            <span hidden={!stop}>入力がないため読み込み停止中</span>
          </Col>
        </Row>
        {masterId && (
          <Row gutter={16} className={paddingBottom}>
            <Form layout="inline">
              <Form.Item label="アプリ入力用キー：">
                <Input value={masterId} readOnly />
              </Form.Item>
              <Form.Item>
                <Button
                  onClick={() => {
                    navigator.clipboard.writeText(masterId);
                    notification['success']({ message: 'Copied!' });
                  }}
                >
                  Copy
                </Button>
              </Form.Item>
            </Form>
          </Row>
        )}
        <Form
          form={form}
          onFinish={send}
          initialValues={{
            name: localStorage.getItem('name') ?? ''
          }}
          validateTrigger="onFinish"
        >
          <Row gutter={16}>
            <Col span={10}>
              <Form.Item name="name" label="名前(空欄可)">
                <Input maxLength={30} disabled={sending} />
              </Form.Item>
            </Col>
            <Col span={10}>
              <Form.Item
                name="url"
                label="URL"
                rules={[{ required: true, type: 'url' }]}
                requiredMark={'optional'}
              >
                <Input
                  placeholder="youtube/niconico/twitter/spotify"
                  disabled={sending}
                />
              </Form.Item>
            </Col>
            <Col span={4}>
              <Form.Item
                name="hidden"
                label="URLを隠す"
                valuePropName="checked"
              >
                <Checkbox />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={10}>
              <Form.Item name="title" label="カスタムタイトル">
                <Input
                  placeholder="未入力の場合はURLから取得"
                  disabled={sending}
                />
              </Form.Item>
            </Col>
            <Col span={14}>
              <Form.Item name="start" label="start/end">
                <TimePicker.RangePicker disabled={sending} order={false} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16} align="middle">
            <Col span={12}>
              <div className={right}>
                <Text>{data?.info.calcTime}</Text>
              </div>
            </Col>
            <Col span={12}>
              <Button type="primary" htmlType="submit" disabled={sending}>
                送信
              </Button>
            </Col>
          </Row>
        </Form>
      </div>
      <DataTable
        data={data}
        roomId={match!.params.roomId!}
        onDelete={() => getList()}
      />
    </Content>
  );
};

export default Room;

const DataTable = (d: {
  data: GetResponse | undefined;
  roomId: string;
  onDelete?: () => void;
}) => {
  const { theme } = useContext(ColorContext);
  const [deleting, setDeleting] = useState(false);
  const [update, setUpdate] = useState(true);
  const c = theme == 'dark' ? currentRowDark : currentRowLight;
  const { useBreakpoint } = Grid;
  const elm = useRef<HTMLDivElement>(null);
  const [height, setHeight] = useState('100vh');
  const [cookies] = useCookies(['uuid']);

  const remove = (idx: number) => {
    setDeleting(true);
    if (!d.data) {
      return;
    }
    const params = new URLSearchParams();
    const index = d.data.data.length - idx - 1;
    params.append('index', `${index}`);
    params.append('uuid', cookies.uuid);
    params.append('room_id', d.roomId);
    axios
      .post(`${ROOT()}/youtube-list/api/youtubelist/remove`, params)
      .then(() => {
        lastUpdateDate = moment(0);
        d.onDelete?.();
      })
      .catch((e: AxiosError) => {
        notification['error']({
          message: 'Error',
          description: e.response?.data.error ?? e.message
        });
      })
      .finally(() => {
        setDeleting(false);
      });
  };

  useEffect(() => {
    if (update) {
      setUpdate(false);
    }
    const resizefunc = () => {
      setUpdate(true);
    };
    window.addEventListener('resize', resizefunc);
    setTimeout(() => {
      if (elm.current) {
        setHeight(
          `calc( 100vh - ${elm.current.getBoundingClientRect().top + 60}px)`
        );
      }
    }, 1000);
    return () => window.removeEventListener('resize', resizefunc);
  }, [update]);
  return (
    <div className={paddingCss} hidden={update}>
      <Table
        ref={elm}
        dataSource={d.data?.data}
        size="small"
        pagination={false}
        scroll={{ y: height }}
        rowClassName={(r: Data, idx) => {
          return r.deleted
            ? deletedRow
            : d.data?.data &&
              d.data?.info.currentIndex === d.data.data.length - idx - 1
            ? c
            : '';
        }}
      >
        <Column
          title="時刻"
          dataIndex="time"
          responsive={['lg']}
          width="150px"
        />
        <Column
          title="送信者"
          dataIndex="username"
          width={useBreakpoint().lg ? '150px' : '75px'}
          render={(text, record: Data, index) => (
            <Space style={{ wordBreak: 'break-all' }}>
              {text}
              {record.removable && (
                <Button
                  size="small"
                  disabled={deleting}
                  danger
                  onClick={() => remove(index)}
                >
                  削除
                </Button>
              )}
            </Space>
          )}
        />
        <Column title="再生時間" dataIndex="length" width="90px" />
        <Column
          title="タイトル"
          dataIndex="title"
          render={(_, r: Data) => (
            <AntdLink href={r.url} target="_blank">
              <Space align="center">
                <div className={headerCss}>{r.icon}</div>
                {r.title}
              </Space>
            </AntdLink>
          )}
        />
      </Table>
    </div>
  );
};
