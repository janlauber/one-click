/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("6jgxtwxj1fdutrh")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "qdwfyqcb",
    "name": "color",
    "type": "select",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "default",
        "dark",
        "red",
        "green",
        "yellow",
        "indigo",
        "purple",
        "pink"
      ]
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("6jgxtwxj1fdutrh")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "qdwfyqcb",
    "name": "color",
    "type": "select",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "red",
        "green",
        "orange",
        "yellow",
        "blue",
        "gray"
      ]
    }
  }))

  return dao.saveCollection(collection)
})
