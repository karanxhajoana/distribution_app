'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface PackResult {
  orderQuantity: number;
  packs: { [key: number]: number };
  totalItems: number;
}

interface PackSizesResponse {
  sizes: number[];
}

export default function Home() {
  const [packSizes, setPackSizes] = useState<number[]>([]);
  const [newSize, setNewSize] = useState('');
  const [quantity, setQuantity] = useState('');
  const [result, setResult] = useState<PackResult | null>(null);
  const [error, setError] = useState('');
  const [editingSize, setEditingSize] = useState<number | null>(null);
  const [editValue, setEditValue] = useState('');

  useEffect(() => {
    fetchPackSizes();
  }, []);

  const fetchPackSizes = async () => {
    try {
      const response = await axios.get<PackSizesResponse>(`${API_URL}/api/pack-sizes`);
      setPackSizes(response.data.sizes);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError('Failed to fetch pack sizes');
    }
  };

  const handleAddSize = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSize) return;

    try {
      await axios.post(`${API_URL}/api/pack-sizes`, {
        size: parseInt(newSize)
      });
      setNewSize('');
      fetchPackSizes();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError('Failed to add pack size');
    }
  };

  const handleDeleteSize = async (size: number) => {
    try {
      await axios.delete(`${API_URL}/api/pack-sizes/${size}`);
      fetchPackSizes();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError('Failed to delete pack size');
    }
  };

  const handleEditSize = async (oldSize: number, newSize: number) => {
    try {
      await axios.put(`${API_URL}/api/pack-sizes`, {
        oldSize,
        newSize
      });
      setEditingSize(null);
      setEditValue('');
      fetchPackSizes();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError('Failed to update pack size');
    }
  };

  const handleCalculate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!quantity) return;

    try {
      const response = await axios.get<PackResult>(`${API_URL}/api/calculate?quantity=${quantity}`);
      setResult(response.data);
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      setError('Failed to calculate pack sizes');
    }
  };

  return (
    <main className="min-h-screen p-8 max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Pack Sizer</h1>
      
      {/* Pack Sizes Section */}
      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Current Pack Sizes</h2>
        <div className="flex flex-wrap gap-2 mb-4">
          {packSizes.map((size) => (
            <div key={size} className="bg-gray-100 px-4 py-2 rounded flex items-center gap-3 shadow-sm">
              {editingSize === size ? (
                <>
                  <input
                    type="number"
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    className="w-24 border p-1 rounded text-lg font-medium text-black"
                    autoFocus
                  />
                  <button
                    onClick={() => handleEditSize(size, parseInt(editValue))}
                    className="text-green-600 hover:text-green-800 text-lg"
                  >
                    ✓
                  </button>
                  <button
                    onClick={() => {
                      setEditingSize(null);
                      setEditValue('');
                    }}
                    className="text-red-600 hover:text-red-800 text-lg"
                  >
                    ✕
                  </button>
                </>
              ) : (
                <>
                  <span className="text-lg font-medium text-black">{size}</span>
                  <div className="flex gap-2">
                    <button
                      onClick={() => {
                        setEditingSize(size);
                        setEditValue(size.toString());
                      }}
                      className="text-blue-600 hover:text-blue-800 text-lg"
                      title="Edit size"
                    >
                      ✎
                    </button>
                    <button
                      onClick={() => handleDeleteSize(size)}
                      className="text-red-600 hover:text-red-800 text-lg"
                      title="Delete size"
                    >
                      ×
                    </button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
        
        <form onSubmit={handleAddSize} className="flex gap-2">
          <input
            type="number"
            value={newSize}
            onChange={(e) => setNewSize(e.target.value)}
            placeholder="New pack size"
            className="border p-2 rounded flex-1 text-lg"
          />
          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 text-lg">
            Add Size
          </button>
        </form>
      </section>

      {/* Calculator Section */}
      <section>
        <h2 className="text-xl font-semibold mb-4">Calculate Packs</h2>
        <form onSubmit={handleCalculate} className="flex gap-2 mb-4">
          <input
            type="number"
            value={quantity}
            onChange={(e) => setQuantity(e.target.value)}
            placeholder="Enter quantity"
            className="border p-2 rounded flex-1"
          />
          <button type="submit" className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
            Calculate
          </button>
        </form>

        {result && (
          <div className="mt-4 bg-white p-4 rounded shadow">
            <h3 className="font-semibold mb-4 text-xl text-black">Pack Distribution:</h3>
            <div className="flex flex-wrap gap-3">
              {Object.entries(result.packs).map(([size, quantity]) => (
                <div key={size} className="bg-green-100 px-4 py-2 rounded shadow-sm">
                  <span className="text-lg font-medium text-black">
                    {quantity} × {size}
                  </span>
                </div>
              ))}
            </div>
            <div className="mt-4 text-sm text-gray-600">
              <p>Order Quantity: {result.orderQuantity}</p>
              <p>Total Items: {result.totalItems}</p>
            </div>
          </div>
        )}
      </section>

      {error && (
        <div className="mt-4 text-red-500 bg-red-50 p-3 rounded">
          {error}
        </div>
      )}
    </main>
  );
}
