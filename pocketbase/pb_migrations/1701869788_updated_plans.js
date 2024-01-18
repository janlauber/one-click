/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "tt9wj0le",
    "name": "manifest",
    "type": "json",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("s7nhljkrmzzu8y6")

  // remove
  collection.schema.removeField("tt9wj0le")

  return dao.saveCollection(collection)
})
