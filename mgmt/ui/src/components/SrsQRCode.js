import React from "react";
import QRCode from "react-qr-code";

export default function SrsQRCode({url}) {
  return (
    <>
      { url ? <QRCode id="qrCode" value={url} size={200} fgColor="#661111" /> : "" }
    </>
  );
}

