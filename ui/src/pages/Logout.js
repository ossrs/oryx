//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import Container from "react-bootstrap/Container";
import React from "react";
import {useNavigate} from "react-router-dom";
import {Token} from "../utils";
import {useTranslation} from "react-i18next";

export default function Logout({onLogout}) {
  const navigate = useNavigate();
  const {t} = useTranslation();

  React.useEffect(() => {
    if (window.confirm(t('nav.logout2'))) {
      Token.remove();
      onLogout && onLogout();
    }

    navigate('/routers-login');
  }, [navigate, onLogout]);

  return <Container>Logout</Container>;
}

