ALTER TABLE nodes 
    ADD COLUMN text_id INT,
    ADD CONSTRAINT node_text 
        FOREIGN KEY (text_id) 
        REFERENCES texts(id);