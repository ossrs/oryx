//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";

export function useIsMobile() {
  const [isMobile, setIsMobile] = React.useState(false);

  function handleWindowSizeChange() {
    setIsMobile(window.innerWidth <= 1100);
  }
  React.useEffect(() => {
    handleWindowSizeChange();
    window.addEventListener('resize', handleWindowSizeChange);
    return () => {
      window.removeEventListener('resize', handleWindowSizeChange);
    }
  }, [setIsMobile]);

  return isMobile;
}
