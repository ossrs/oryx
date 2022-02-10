import React from "react";
import {Button, Spinner, OverlayTrigger, Popover} from "react-bootstrap";
import {useSearchParams} from "react-router-dom";

export default function PopoverConfirmButton({onClick, availableRelease, upgrading, progress, operator, text, children}) {
  const [startUpgrade, setStartUpgrade] = React.useState();
  const [disabled, setDisabled] = React.useState(true);
  const [searchParams] = useSearchParams();

  React.useEffect(() => {
    const allowForceUpgrade = searchParams.get('allow-force') === 'true';
    setDisabled(upgrading || (!availableRelease && !allowForceUpgrade));
  }, [availableRelease, upgrading]);

  React.useEffect(() => {
    const allowForceUpgrade = searchParams.get('allow-force') === 'true';
    console.log(`?allow-force=true|false, current=${allowForceUpgrade}, Whether allow force to upgrade, even it's the latest version`);
  }, []);

  const onHandleClick = () => {
    setStartUpgrade(false);
    onClick && onClick();
  };

  const filterByOperator = (msg) => {
    return operator ? msg.replaceAll('升级', operator) : msg;
  };

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
              disabled={disabled}
              onClick={!disabled ? onHandleClick : null}
            >
              确认
            </Button>
          </div>
          <div className="col-12">
            <Button
              variant="primary"
              onClick={() => setStartUpgrade(false)}
            >
              取消
            </Button>
          </div>
        </div>
      </Popover.Body>
    </Popover>
  );

  return (
    <div className='row row-cols-lg-auto g-3 align-items-center' style={{display: 'inline-block'}}>
      <div className="col-12">
        <OverlayTrigger trigger="click" placement="right" overlay={popover} show={startUpgrade}>
          <Button variant="primary" onClick={() => setStartUpgrade(!startUpgrade)} disabled={disabled}>
            {upgrading ? filterByOperator('升级中...') : text}
          </Button>
        </OverlayTrigger> &nbsp;
        {upgrading && <Spinner animation="border" variant="success" style={{verticalAlign: 'middle'}} />} &nbsp;
        {upgrading && progress}
      </div>
    </div>
  );
}

