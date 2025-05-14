import React from 'react';
import { stores } from '../stores';

interface StoreSelectorProps {
    selectedStoreId: string;
    onStoreChange: (event: React.ChangeEvent<HTMLSelectElement>) => void;
}

function StoreSelector({ selectedStoreId, onStoreChange }: StoreSelectorProps) {
    return (
        <select
            id="store-select"
            value={selectedStoreId}
            onChange={onStoreChange}
            className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        >
            {stores.map(store => (
                <option key={store.id} value={store.id}>
                    {store.name}
                </option>
            ))}
        </select>
    );
}

export default StoreSelector;
