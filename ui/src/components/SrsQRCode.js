//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";
import QRCode from "react-qr-code";

export default function SrsQRCode({url}) {
  if (!url) return <></>;
  return (
    <QRCode value={url} data-testid='qrCode' size={200} fgColor="#661111" />
  );
}

