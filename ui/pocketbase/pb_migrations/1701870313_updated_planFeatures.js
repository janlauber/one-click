/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "rikp7sjv",
    "name": "options",
    "type": "json",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("m9zfrl1ndxsejxd")

  // remove
  collection.schema.removeField("rikp7sjv")

  return dao.saveCollection(collection)
})
