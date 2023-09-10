//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: AGPL-3.0-or-later
//
import React from "react";
import {Button, OverlayTrigger, Popover} from "react-bootstrap";
import {useSrsLanguage} from "./LanguageSwitch";

export default function PopoverConfirm({onClick, trigger, children, placement}) {
  const [startUpgrade, setStartUpgrade] = React.useState();
  const language = useSrsLanguage();

  const onHandleClick = React.useCallback(() => {
    setStartUpgrade(false);
    onClick && onClick();
  }, [onClick]);

  const popover = (
    <Popover id="popover-basic">
      <Popover.Header as="h3">Confirm</Popover.Header>
      <Popover.Body>
        <div>
          {children}
        </div>
        <div className='row row-cols-lg-auto g-3 align-items-center'>
          <div className="col-12">
            <Button
              variant="danger"
              onClick={onHandleClick}
            >
              {language === 'zh' ? '确认' : 'Continue'}
            </Button>
          </div>
          <div className="col-12">
            <Button
              variant="primary"
              onClick={() => setStartUpgrade(false)}
            >
              {language === 'zh' ? '取消' : 'Abort'}
            </Button>
          </div>
        </div>
      </Popover.Body>
    </Popover>
  );

  return (
    <OverlayTrigger trigger="click" placement={placement || 'right'} overlay={popover} show={startUpgrade}>
      <span onClick={(e) => {
        e.preventDefault();
        setStartUpgrade(!startUpgrade);
      }}>
        {trigger}
      </span>
    </OverlayTrigger>
  );
}

