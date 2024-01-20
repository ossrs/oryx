import React from "react";
import {Alert} from "react-bootstrap";

export function AITalkErrorLogPanel({errorLogs, removeErrorLog}) {
  return (
    <React.Fragment>
      {errorLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeErrorLog(log)} variant='danger' dismissible>
            <Alert.Heading>Error!</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
    </React.Fragment>
  );
}

export function AITalkTipLogPanel({tipLogs, removeTipLog}) {
  return (
    <React.Fragment>
      {tipLogs.map((log) => {
        return (
          <Alert key={log.id} onClose={() => removeTipLog(log)} variant='success' dismissible>
            <Alert.Heading>{log.title}</Alert.Heading>
            <p>{log.msg}</p>
          </Alert>
        );
      })}
    </React.Fragment>
  );
}
