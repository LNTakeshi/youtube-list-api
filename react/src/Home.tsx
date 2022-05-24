import { BgColorsOutlined } from '@ant-design/icons';
import { css } from '@linaria/core';
import { Button, Input, Row, Space } from 'antd';
import Form from 'antd/lib/form/Form';
import FormItem from 'antd/lib/form/FormItem';
import { Content } from 'antd/lib/layout/layout';
import Title from 'antd/lib/typography/Title';
import React, { useContext } from 'react';
import { isChrome, isChromium, isEdge } from 'react-device-detect';
import { NavigateFunction, useNavigate } from 'react-router-dom';
import { ColorContext, ForceReloadContext } from './App';

const contentCss = css`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
`;
const padding = css`
  padding: 8px;
`;

export const useChangeColor = () => {
  const { theme, setTheme } = useContext(ColorContext);
  const useForceReload = useContext(ForceReloadContext);

  const changeTheme = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
    if (isChrome || isChromium || isEdge) {
      console.log('chrome');
      window.location.reload();
      useForceReload.setLoad(false);
    }
  };
  return changeTheme;
};

export const ChangeColorButton = (props: { onClick: () => void }) => {
  return (
    <Button onClick={props.onClick}>
      <BgColorsOutlined />
      Color
    </Button>
  );
};

const Home = () => {
  const navigate: NavigateFunction = useNavigate();
  const changeColor = useChangeColor();

  const onFinish = (v: { idform: string }) => {
    if (v.idform) {
      navigate(`${process.env.PUBLIC_URL}/room/${v.idform}`);
    }
  };
  return (
    <>
      <Content className={contentCss}>
        <Title>Youtubeの動画を流すアプリ</Title>
        <Form layout="inline" method="get" onFinish={onFinish}>
          <FormItem name="idform" label="部屋ID">
            <Input />
          </FormItem>
          <FormItem>
            <Button type="primary" htmlType="submit">
              入室
            </Button>
          </FormItem>
        </Form>

        <Row gutter={16} className={padding}>
          <Space>
            <ChangeColorButton onClick={changeColor} />

            <a href="https://drive.google.com/file/d/1-1Rg0ZuO2Ro3bEizkoGxeTrkCfIyG_A5/view?usp=sharing">
              アプリのダウンロード
            </a>
          </Space>
        </Row>
        <Row gutter={16} className={padding}>
          <Space direction="vertical" align="center">
            <div>2.0.0: APIの最適化、他</div>
            <div>
              2.1.3: 色々かわった Vcyncを有効無効切り替えられるように
              Spotifyの画像取得処理の修正
            </div>
          </Space>
        </Row>
      </Content>
    </>
  );
};

export default Home;
