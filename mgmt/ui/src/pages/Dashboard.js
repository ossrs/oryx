import React from "react";
import {Container, Tabs, Tab, Accordion} from "react-bootstrap";
import moment from "moment";
import axios from 'axios';
import {XAxis, Tooltip, CartesianGrid, AreaChart, YAxis, Area} from "recharts";
import querystring from "querystring";

export default function Dashboard() {
  const [data, setData] = React.useState();
  const [promQL, setPromQL] = React.useState();

  React.useEffect(() => {
    const query = '(1 - min by(mode) (rate(node_cpu_seconds_total{mode="idle"}[10s]))) * 100';
    const queryEscaped = querystring.stringify({
      "g0.expr": query,
      "g0.tab": 0,
      "g0.stacked": 0,
      "g0.show_exemplars": 0,
      "g0.range_input": "1h",
    });
    setPromQL(`/prometheus/graph?${queryEscaped}`);

    // See https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
    axios.get(`/prometheus/api/v1/query_range`, {
      params: {
        query: query,
        start: moment().subtract(60, 'minutes').valueOf() / 1000.0,
        end: moment().valueOf() / 1000.0,
        step: 30,
      },
    }).then(res => {
      const matrix = res.data.data.result;
      if (matrix && matrix.length && matrix[0].values) {
        const samples = matrix[0].values.map(e => ({
          time: moment.unix(e[0]).format('HH:mm'),
          cpu: Number(e[1]).toFixed(1),
        }));
        console.log('cpu samples is', samples);
        setData(samples);
      }
      console.log(`Status: Query ok, matrix=${matrix.length}`);
    }).catch(e => {
      const err = e.response.data;
      alert(`服务器错误，${err.code}: ${err.data.message}`);
    });
  }, []);

  return (
    <>
      <p></p>
      <Container>
        <Tabs defaultActiveKey="sys" id="uncontrolled-tab-example" className="mb-3">
          <Tab eventKey="sys" title="系统状态">
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>CPU</Accordion.Header>
                <Accordion.Body>
                  <AreaChart
                    width={1270}
                    height={400}
                    data={data}
                    margin={{ top: 10, right: 30, left: 0, bottom: 10 }}
                  >
                    <defs>
                      <linearGradient id="colorCpu" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#82ca9d" stopOpacity={0}/>
                        <stop offset="95%" stopColor="#82ca9d" stopOpacity={0.8}/>
                      </linearGradient>
                    </defs>
                    <XAxis dataKey="time" />
                    <YAxis type='number' domain={[0, 100]} />
                    <CartesianGrid strokeDasharray="3 3" />
                    <Tooltip />
                    <Area type="monotone" dataKey="cpu" stroke="#82ca9d" fillOpacity={1} fill="url(#colorCpu)" />
                  </AreaChart>
                  <a href={promQL} target='_blank' rel='noreferrer'>Show in Prometheus</a>
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Tab>
        </Tabs>
      </Container>
    </>
  );
}

