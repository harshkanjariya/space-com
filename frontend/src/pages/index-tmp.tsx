import { useState } from 'react';
import {post} from "../common/api/axios-functions.ts";

export function Index() {
  const [object1Message, setObject1Message] = useState('');
  const [object2Message, setObject2Message] = useState('');
  const [object3Message, setObject3Message] = useState('');

  // State to store results
  const [conversionResult, setConversionResult] = useState('');

  // Handler for broadcasting messages
  const handleBroadcast = async (format: string, message: string) => {
    try {
      const response = await post('/converter/encode', {
        outputFormat: format,
        message: message,
      });
      setConversionResult(response.data);
    } catch (error) {
      console.log(error);
      setConversionResult('Error broadcasting message');
    }
  };

  return (
    <div>
      <h1>Object Messaging</h1>

      <div>
        <h2>Object 1 (Format: AOS)</h2>
        <textarea
          placeholder="Enter message for Object 1"
          value={object1Message}
          onChange={(e) => setObject1Message(e.target.value)}
        />
        <button onClick={() => handleBroadcast('aos', object1Message)}>Broadcast Object 1</button>
      </div>

      <div>
        <h2>Object 2 (Format: PUS_TM)</h2>
        <textarea
          placeholder="Enter message for Object 2"
          value={object2Message}
          onChange={(e) => setObject2Message(e.target.value)}
        />
        <button onClick={() => handleBroadcast('pus_tm', object2Message)}>Broadcast Object 2</button>
      </div>

      <div>
        <h2>Object 3 (Format: PUS_TC)</h2>
        <textarea
          placeholder="Enter message for Object 3"
          value={object3Message}
          onChange={(e) => setObject3Message(e.target.value)}
        />
        <button onClick={() => handleBroadcast('pus_tc', object3Message)}>Broadcast Object 3</button>
      </div>

      <div>
        <h3>Conversion Result</h3>
        <pre>{conversionResult}</pre>
      </div>
    </div>
  );
}
