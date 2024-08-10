import {useEffect, useState} from 'react';
import {Button, Input, Modal} from 'antd';
import 'antd/dist/reset.css';
import {post} from '../common/api/axios-functions.ts';

export function Index() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [currentSatellite, setCurrentSatellite] = useState<string>('');
  const [inputValue, setInputValue] = useState('');
  const [satelliteData, setSatelliteData] = useState<Record<string, string>>({});
  const [satelliteMessages, setSatelliteMessages] = useState<Record<string, string>>({});
  const [tooltipDisplay, setTooltipDisplay] = useState("none");
  const [tooltipPosition, setTooltipPosition] = useState({top: -100, left: -100});
  const [tooltipContent, setTooltipContent] = useState("");

  const satelliteToOutputFormat: Record<"satellite1" | "satellite2" | "satellite3" | "common", string> = {
    satellite1: 'aos',
    satellite2: 'pus_tm',
    satellite3: 'pus_tc',
    common: "ccsds",
  };

  async function convertFormat(inputFormat: string, outputFormat: string, data: string) {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const result = await post("/converter/convert", {
      inputFormat,
      outputFormat,
      data,
    });
    return result.split("\n");
  }

  async function decode(inputFormat: string, data: string) {
    return await post("/converter/decode", {
      inputFormat,
      data,
    });
  }

  useEffect(() => {
    let payload: any;
    for (const type of Object.keys(satelliteToOutputFormat)) {
      if (satelliteData[type] && type != "common") {
        payload = {
          data: satelliteData[type],
          format: (satelliteToOutputFormat as any)[type],
        };
        break;
      }
    }
    if (payload && !satelliteData.common) {
      convertFormat(payload.format, satelliteToOutputFormat.common, payload.data).then((response) => {
        setSatelliteData({
          ...satelliteData,
          "common": response[response.length - 2],
        });
      });
    }
    if (satelliteData.common) {
      for (const type of Object.keys(satelliteToOutputFormat)) {
        if (type == "common") continue;
        if (!satelliteData[type]) {
          convertFormat(
            satelliteToOutputFormat.common,
            (satelliteToOutputFormat as any)[type],
            satelliteData.common,
          ).then((response) => {
            setSatelliteData({
              ...satelliteData,
              [type]: response[response.length - 2],
            })
          });
        } else if (!satelliteMessages[type]) {
          decode((satelliteToOutputFormat as any)[type], satelliteData[type]).then((result) => {
            setSatelliteMessages({
              ...satelliteMessages,
              [type]: result,
            })
          });
        }
      }
    }
  }, [satelliteData]);

  const showModal = (satellite: "satellite1" | "satellite2" | "satellite3") => {
    setCurrentSatellite(satellite);
    setIsModalVisible(true);
  };

  const handleOk = async () => {
    try {
      const encoded = await post('/converter/encode', {
        outputFormat: (satelliteToOutputFormat as any)[currentSatellite],
        message: inputValue,
      });
      const lines = encoded.split("\n");
      setSatelliteData({
        ...satelliteData,
        [currentSatellite]: lines[lines.length - 2],
      });
      setSatelliteMessages({
        [currentSatellite]: "\u0000",
      })

      // Close the modal and reset states
      setIsModalVisible(false);
      setInputValue('');
    } catch (error) {
      console.error('Error calling API:', error);
    }
  };

  const handleCancel = () => {
    setIsModalVisible(false);
    setInputValue('');
  };

  const handleMouseEnter = (satellite: "satellite1" | "satellite2" | "satellite3", event: React.MouseEvent) => {
    const data = satelliteData[satellite];
    const message = satelliteMessages[satellite];

    const content = `protocol: ${(satelliteToOutputFormat as any)[satellite]} ${data ? "/ data: " + data : ""} ${message && message !== "\u0000" ? `\nmessage received: ${message}` : ""}`;

    setTooltipContent(content);
    const obj = (event.target as any).getBoundingClientRect();
    const {left, bottom} = obj;
    setTooltipPosition({
      top: bottom + 10,
      left: left + 10,
    });
    setTooltipDisplay("block");
  };

  const handleMouseLeave = () => {
    setTooltipDisplay("none");
  };

  return (<>
      <div className="container">
        <Button
          onClick={() => {
            setSatelliteMessages({});
            setSatelliteData({});
          }}
          className="fixed"
          style={{top: -150, left: 50}}>
          Reset
        </Button>
        <div className="converter" id="converter1"></div>
        <div className="converter" id="converter2"></div>
        <div className="converter" id="converter3"></div>

        <div
          className="satellite"
          id="satellite1"
          onClick={() => showModal('satellite1')}
        ></div>
        <span
          onMouseEnter={(event) => handleMouseEnter('satellite1', event)}
          onMouseLeave={handleMouseLeave}
          className="label label1"
          onClick={() => showModal('satellite1')}
        >
        protocol: aos
          {satelliteData?.satellite1 && <div>{"data: " + satelliteData?.satellite1}</div>}
          {satelliteMessages?.satellite1 && satelliteMessages?.satellite1 !== "\u0000" && <>
              message received: {satelliteMessages?.satellite1}
          </>}
      </span>

        <div
          className="satellite"
          id="satellite2"
          onClick={() => showModal('satellite2')}
        ></div>
        <span
          onMouseEnter={(event) => handleMouseEnter('satellite2', event)}
          onMouseLeave={handleMouseLeave}
          className="label label2"
          onClick={() => showModal('satellite2')}
        >
        protocol: pus-tm
          {satelliteData?.satellite2 && <div>{"data: " + satelliteData?.satellite2}</div>}
          {satelliteMessages?.satellite2 && satelliteMessages?.satellite2 !== "\u0000" && <>
              message received: {satelliteMessages?.satellite2}
          </>}
      </span>

        <div
          className="satellite"
          id="satellite3"
          onClick={() => showModal('satellite3')}
        ></div>
        <span
          onMouseEnter={(event) => handleMouseEnter('satellite3', event)}
          onMouseLeave={handleMouseLeave}
          className="label label3"
          onClick={() => showModal('satellite3')}
        >
        protocol: pus-tc
          {satelliteData?.satellite3 && <div>{"data: " + satelliteData?.satellite3}</div>}
          {satelliteMessages?.satellite3 && satelliteMessages?.satellite3 !== "\u0000" && <>
              message received: {satelliteMessages?.satellite3}
          </>}
      </span>

        {satelliteData.common && <span className="label label4">{satelliteData?.common}</span>}

        <div className="line" id="line1"></div>
        <div className="line" id="line2"></div>
        <div className="line" id="line3"></div>

        <div className="line" id="line4"></div>
        <div className="line" id="line5"></div>
        <div className="line" id="line6"></div>

        <Modal
          title={`Input for ${currentSatellite}`}
          open={isModalVisible}
          onOk={handleOk}
          onCancel={handleCancel}
        >
          <Input
            placeholder="Enter data"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
          />
        </Modal>

      </div>
      <div style={{...tooltipPosition, display: tooltipDisplay}} className="tooltip">
        {tooltipContent.split("\n").map((line, index) => (
          <div key={index}>{line}</div>
        ))}
      </div>
    </>
  );
}
