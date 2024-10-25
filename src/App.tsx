import { Actor } from './Actor';
import './App.css';

import { Stage } from '@pixi/react';

const App = () => {
  return (
    <Stage width={window.innerWidth} height={window.innerHeight} options={{ background: 0x1099bb }}>
      <Actor />
      <Actor />
      <Actor />
    </Stage>
  );
};

export default App;