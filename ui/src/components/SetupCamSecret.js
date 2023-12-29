//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Button, Form} from "react-bootstrap";
import {Token} from "../utils";
import axios from "axios";
import {useTranslation} from "react-i18next";
import {useErrorHandler} from "react-error-boundary";

export default function SetupCamSecret({children}) {
  const [secretId, setSecretId] = React.useState();
  const [secretKey, setSecretKey] = React.useState();
  const handleError = useErrorHandler();
  const {t} = useTranslation();

  const updateTencentSecret = React.useCallback((e) => {
    e.preventDefault();

    axios.post('/terraform/v1/tencent/cam/secret', {
      secretId, secretKey,
    }, {
      headers: Token.loadBearerHeader(),
    }).then(res => {
      alert(t('tencent.secretOk'));
    }).catch(handleError);
  }, [handleError, secretId, secretKey, t]);

  return (<>
    <Form>
      <Form.Group className="mb-3">
        <Form.Label>SecretId</Form.Label>
        <Form.Text> * {t('tencent.secretIdTip')}, <a href='https://console.cloud.tencent.com/cam/capi' target='_blank' rel='noreferrer'>{t('tencent.secretGet')}</a></Form.Text>
        <Form.Control as="input" rows={2} defaultValue={secretId} onChange={(e) => setSecretId(e.target.value)} />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>SecretKey</Form.Label>
        <Form.Text> * {t('tencent.secretKeyTip')}, <a href='https://console.cloud.tencent.com/cam/capi' target='_blank' rel='noreferrer'>{t('tencent.secretGet')}</a></Form.Text>
        <Form.Control as="input" type='password' rows={2} defaultValue={secretKey} onChange={(e) => setSecretKey(e.target.value)} />
      </Form.Group>
      <Button variant="primary" type="submit" onClick={(e) => updateTencentSecret(e)}>
        {t('tencent.secretSubmit')}
      </Button>
      <Form.Text> * {t('tencent.secretSubmitTip')}</Form.Text> &nbsp;
      {children}
    </Form>
  </>);
}

