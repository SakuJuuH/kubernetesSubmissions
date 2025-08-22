import { useState, useEffect } from 'react'
import axios from 'axios'
import './App.css'

interface ImageInfo {
  path: string;
  cached_at: string;  
}

function App() {
  const [imageInfo, setImageInfo] = useState<ImageInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [todo, setTodo] = useState<string>('');

  let [todos, setTodos] = useState<string[]>([
      "Buy groceries", "Walk the dog", "Read a book"
  ]);

  const backendUrl = import.meta.env.PROD 
    ? '/api'
    : 'http://localhost:3000/api';

  useEffect(() => {
      const fetchImageInfo = async () => {
      try {
        setLoading(true)
        setError(null);
    
        const response = await axios.get<ImageInfo>(`${backendUrl}/image`);
        
        if (!response.data || !response.data.path) {
          throw new Error('Invalid image data received');
        }

        setImageInfo(response.data);
        console.log(response.data);

      } catch (error) {
        console.error('Error fetching image:', error);
        setError('Failed to fetch image');
      } finally {
        setLoading(false);
      }
    };

    fetchImageInfo().then();
  }, [backendUrl]);

  const handleShutdown = async () => {
    try {
      await axios.post(`${backendUrl}/shutdown`);
      alert('Server shutdown initiated!');
    } catch (error) {
      console.error('Error shutting down server:', error);
    }
  };

  const handleAddTodo = async () => {
      if (!todo.trim()) {
        alert('Please enter a todo item.');
        return;
      }

      if (todo.length > 140) {
        alert('Todo item cannot exceed 140 characters.');
        return;
      }
      todos.push(todo);
      setTodos([...todos]);
      setTodo('');
  }

  return (
    <>
      <div>
        <h1>The Todo App</h1>
        {loading && <p>Loading image...</p>}
        {error && <p style={{color: 'red'}}>{error}</p>}
        
        {imageInfo && !loading && (
          <div>
            <p>Image cached at: {new Date(imageInfo.cached_at).toLocaleString()}</p>
            <img
                src={`${backendUrl}${imageInfo.path}`}
                alt="Random image from Picsum"
                style={{maxWidth: '100%', height: 'auto'}}
            />
          </div>
        )}
          <div className="todo" >
              <input
                  className="input"
              type="text"
              placeholder=""
              value={todo}
              onChange={(e) => setTodo(e.target.value)}
          />
              <button className="button" onClick={handleAddTodo}>
                Add Todo
              </button>
          </div>

          <div className="todo-list">
              {
                  todos ? (todos.map((item, index) => (
                      <div key={index} className="todo-item">
                        <span>{index + 1}. {item}</span>
                      </div>
                  ))) : "No todos available"
              }
          </div>
        <div style={{marginTop: '20px'}}>
          <button onClick={handleShutdown}>
            Shutdown Server (for testing)
          </button>
        </div>
      </div>
    </>
  );
}

export default App;